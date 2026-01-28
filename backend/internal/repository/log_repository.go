package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"log-management-backend/internal/models"

	"github.com/google/uuid"
)

type LogRepository struct {
	DB *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{DB: db}
}

// CreateBatch บันทึก Log ทีละหลายๆ ตัว (เพื่อ performance)
func (r *LogRepository) CreateBatch(ctx context.Context, logs []models.LogEntry) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare Statement เพื่อความเร็วและความปลอดภัย
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO logs (id, tenant_id, timestamp, source, event_type, severity, body, raw_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, log := range logs {
		// แปลง Body กลับเป็น JSONB string
		bodyBytes, _ := json.Marshal(log.Body)

		// สร้าง UUID ถ้ายังไม่มี
		if log.ID == uuid.Nil {
			log.ID = uuid.New()
		}

		// Execute Insert
		_, err := stmt.ExecContext(ctx,
			log.ID,
			log.TenantID,
			log.Timestamp,
			log.Source,
			log.EventType,
			log.Severity,
			bodyBytes,
			log.RawMessage,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetLogs - แก้ไขให้คืนค่า []map[string]interface{} แทน
func (r *LogRepository) GetLogs(ctx context.Context, tenantID string, query string, limit int) ([]map[string]interface{}, error) {
	// SQL Query พื้นฐาน: เรียงตามเวลาล่าสุด
	var stmt string
	var args []interface{}

	if tenantID == "" {
		// Admin ดูทั้งหมด (ไม่ filter tenant)
		stmt = `
			SELECT timestamp, tenant_id, source, event_type, severity, raw_message
			FROM logs
			WHERE ($1 = '' OR raw_message ILIKE '%' || $1 || '%')
			ORDER BY timestamp DESC
			LIMIT $2`
		args = []interface{}{query, limit}
	} else {
		// Filter ตาม tenant
		stmt = `
			SELECT timestamp, tenant_id, source, event_type, severity, raw_message
			FROM logs
			WHERE tenant_id = $1
			AND ($2 = '' OR raw_message ILIKE '%' || $2 || '%')
			ORDER BY timestamp DESC
			LIMIT $3`
		args = []interface{}{tenantID, query, limit}
	}

	rows, err := r.DB.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var timestamp, tenantIDVal, source, eventType, rawMessage string
		var severity int

		if err := rows.Scan(&timestamp, &tenantIDVal, &source, &eventType, &severity, &rawMessage); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		logs = append(logs, map[string]interface{}{
			"@timestamp": timestamp,
			"tenant_id":  tenantIDVal,
			"source":     source,
			"event_type": eventType,
			"severity":   severity,
			"raw":        rawMessage,
		})
	}

	return logs, nil
}

// DeleteOldLogs ลบ logs ที่เก่ากว่าจำนวนวันที่กำหนด
func (r *LogRepository) DeleteOldLogs(ctx context.Context, retentionDays int) (int64, error) {
	query := `DELETE FROM logs WHERE timestamp < NOW() - INTERVAL '1 day' * $1`
	result, err := r.DB.ExecContext(ctx, query, retentionDays)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
