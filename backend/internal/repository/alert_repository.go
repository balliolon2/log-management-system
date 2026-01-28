package repository

import (
	"context"
	"database/sql"
	"log"
	"time"

	"log-management-backend/internal/models"
)

type AlertRepository struct {
	db *sql.DB
}

func NewAlertRepository(db *sql.DB) *AlertRepository {
	return &AlertRepository{db: db}
}

// CreateAlert - บันทึก alert ลง database
func (r *AlertRepository) CreateAlert(ctx context.Context, alert models.AlertHistory) error {
	query := `
		INSERT INTO alert_history (tenant_id, rule_id, rule_name, triggered_at, details, severity, acknowledged)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query,
		alert.TenantID,
		alert.RuleID,
		alert.RuleName,
		alert.TriggeredAt,
		alert.Details,
		alert.Severity,
		alert.Acknowledged,
	)
	if err != nil {
		log.Printf("Failed to create alert: %v", err)
		return err
	}
	return nil
}

// GetAlertHistory - ดึงประวัติ alerts (รองรับ pagination)
func (r *AlertRepository) GetAlertHistory(ctx context.Context, tenantID string, limit int) ([]models.AlertHistory, error) {
	query := `
		SELECT id, tenant_id, rule_id, rule_name, triggered_at, details, severity, acknowledged
		FROM alert_history
		WHERE tenant_id = $1
		ORDER BY triggered_at DESC
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, tenantID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alerts []models.AlertHistory
	for rows.Next() {
		var a models.AlertHistory
		err := rows.Scan(
			&a.ID,
			&a.TenantID,
			&a.RuleID,
			&a.RuleName,
			&a.TriggeredAt,
			&a.Details,
			&a.Severity,
			&a.Acknowledged,
		)
		if err != nil {
			log.Printf("Failed to scan alert: %v", err)
			continue
		}
		alerts = append(alerts, a)
	}
	return alerts, nil
}

// GetAlertRules - ดึงกฎทั้งหมดที่ enabled
func (r *AlertRepository) GetAlertRules(ctx context.Context) ([]models.AlertRule, error) {
	query := `
		SELECT id, rule_name, description, event_type, condition_field,
		       condition_operator, threshold, time_window_seconds, severity, enabled, created_at
		FROM alert_rules
		WHERE enabled = true
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.AlertRule
	for rows.Next() {
		var rule models.AlertRule
		err := rows.Scan(
			&rule.ID,
			&rule.RuleName,
			&rule.Description,
			&rule.EventType,
			&rule.ConditionField,
			&rule.ConditionOperator,
			&rule.Threshold,
			&rule.TimeWindowSeconds,
			&rule.Severity,
			&rule.Enabled,
			&rule.CreatedAt,
		)
		if err != nil {
			log.Printf("Failed to scan rule: %v", err)
			continue
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

// CountRecentEvents - นับจำนวน events ในช่วงเวลาที่กำหนด
// ใช้สำหรับตรวจสอบว่าเกิน threshold หรือไม่
func (r *AlertRepository) CountRecentEvents(ctx context.Context, tenantID, eventType, conditionField, conditionValue string, timeWindowSeconds int) (int, error) {
	// คำนวณเวลาย้อนหลัง
	since := time.Now().Add(-time.Duration(timeWindowSeconds) * time.Second)

	query := `
		SELECT COUNT(*)
		FROM logs
		WHERE tenant_id = $1
		  AND event_type = $2
		  AND body ->> $3 = $4
		  AND timestamp >= $5
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, tenantID, eventType, conditionField, conditionValue, since).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
