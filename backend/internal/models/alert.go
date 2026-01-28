package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// AlertRule - โครงสร้างสำหรับกฎการแจ้งเตือน
type AlertRule struct {
	ID                  string    `json:"id"`
	RuleName            string    `json:"rule_name"`
	Description         string    `json:"description"`
	EventType           string    `json:"event_type"`
	ConditionField      string    `json:"condition_field"`
	ConditionOperator   string    `json:"condition_operator"` // "count_gt", "count_gte", "equals"
	Threshold           int       `json:"threshold"`
	TimeWindowSeconds   int       `json:"time_window_seconds"`
	Severity            int       `json:"severity"`
	Enabled             bool      `json:"enabled"`
	CreatedAt           time.Time `json:"created_at"`
}

// AlertHistory - โครงสร้างสำหรับประวัติ alerts
type AlertHistory struct {
	ID            string          `json:"id"`
	TenantID      string          `json:"tenant_id"`
	RuleID        string          `json:"rule_id"`
	RuleName      string          `json:"rule_name"`
	TriggeredAt   time.Time       `json:"triggered_at"`
	Details       AlertDetails    `json:"details"`
	Severity      int             `json:"severity"`
	Acknowledged  bool            `json:"acknowledged"`
}

// AlertDetails - รายละเอียดของ Alert (เก็บใน JSONB)
type AlertDetails map[string]interface{}

// Value - สำหรับ database/sql driver (แปลง struct -> JSON)
func (a AlertDetails) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan - สำหรับ database/sql driver (แปลง JSON -> struct)
func (a *AlertDetails) Scan(value interface{}) error {
	if value == nil {
		*a = make(AlertDetails)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, a)
}
