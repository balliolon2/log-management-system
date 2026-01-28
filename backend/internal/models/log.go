package models

import (
	"time"
	"github.com/google/uuid"
	"encoding/json"
)

// LogSchema
type LogEntry struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	TenantID   string          `json:"tenant" db:"tenant_id"`
	Timestamp  time.Time       `json:"@timestamp" db:"timestamp"`
	Source     string          `json:"source" db:"source"`
	EventType  string          `json:"event_type" db:"event_type"`
	Severity   int             `json:"severity" db:"severity"`

	// ใช้ json.RawMessage หรือ map[string]interface{} สำหรับ JSONB
	Body       json.RawMessage `json:"body" db:"body"`
	RawMessage string          `json:"raw" db:"raw_message"`
}

// โครงสร้างสำหรับรับ Request (API Ingest)
type IngestRequest struct {
    TenantID string                 `json:"tenant"`
    Source   string                 `json:"source"`
    Logs     []map[string]interface{} `json:"logs"` // รองรับ Batch ingestion
}
