package alert

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"log-management-backend/internal/models"
	"log-management-backend/internal/repository"
)

type AlertEngine struct {
	alertRepo   *repository.AlertRepository
	logRepo     *repository.LogRepository
	ruleManager *RuleManager
}

func NewAlertEngine(alertRepo *repository.AlertRepository, logRepo *repository.LogRepository) *AlertEngine {
	return &AlertEngine{
		alertRepo:   alertRepo,
		logRepo:     logRepo,
		ruleManager: NewRuleManager(alertRepo),
	}
}

// Start - เริ่มต้น Alert Engine
func (ae *AlertEngine) Start(ctx context.Context) error {
	// โหลดกฎจาก database
	if err := ae.ruleManager.LoadRules(ctx); err != nil {
		return fmt.Errorf("failed to load rules: %v", err)
	}
	log.Println("🚨 Alert Engine started")
	return nil
}

// CheckLog - ตรวจสอบ log ที่เข้ามาใหม่ว่าตรงกับ rule ไหนไหม (Real-time)
func (ae *AlertEngine) CheckLog(ctx context.Context, logEntry models.LogEntry) {
	// หา rules ที่ตรงกับ event_type ของ log
	matchedRules := ae.ruleManager.GetRuleByEventType(logEntry.EventType)

	for _, rule := range matchedRules {
		// ตรวจสอบ condition
		if ae.shouldTriggerAlert(ctx, logEntry, rule) {
			ae.triggerAlert(ctx, logEntry, rule)
		}
	}
}

// shouldTriggerAlert - เช็คว่าควร trigger alert หรือไม่
func (ae *AlertEngine) shouldTriggerAlert(ctx context.Context, logEntry models.LogEntry, rule models.AlertRule) bool {
	// แปลง Body ([]byte) เป็น map
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(logEntry.Body, &bodyMap); err != nil {
		log.Printf("Failed to unmarshal body: %v", err)
		return false
	}

	// ดึงค่า field ที่ต้องเช็ค (เช่น src_ip, host)
	conditionValueRaw, ok := bodyMap[rule.ConditionField]
	if !ok {
		return false
	}

	// แปลงเป็น string
	conditionValue, ok := conditionValueRaw.(string)
	if !ok {
		// ถ้าไม่ใช่ string ให้ลอง convert
		conditionValue = fmt.Sprintf("%v", conditionValueRaw)
	}

	// นับจำนวน events ในช่วงเวลาที่กำหนด
	count, err := ae.alertRepo.CountRecentEvents(
		ctx,
		logEntry.TenantID,
		rule.EventType,
		rule.ConditionField,
		conditionValue,
		rule.TimeWindowSeconds,
	)
	if err != nil {
		log.Printf("❌ Failed to count events: %v", err)
		return false
	}

	// เช็ค threshold
	switch rule.ConditionOperator {
	case "count_gt": // มากกว่า
		return count > rule.Threshold
	case "count_gte": // มากกว่าหรือเท่ากับ
		return count >= rule.Threshold
	case "equals": // เท่ากับ
		return count == rule.Threshold
	default:
		return false
	}
}

// triggerAlert - สร้าง alert ลง database
func (ae *AlertEngine) triggerAlert(ctx context.Context, logEntry models.LogEntry, rule models.AlertRule) {
	// แปลง Body เป็น map
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(logEntry.Body, &bodyMap); err != nil {
		log.Printf("Failed to unmarshal body: %v", err)
		return
	}

	// ดึงค่า field ที่เกี่ยวข้อง
	conditionValueRaw, _ := bodyMap[rule.ConditionField]
	conditionValue := fmt.Sprintf("%v", conditionValueRaw)

	// นับจำนวนครั้งที่เกิดขึ้น (เพื่อแสดงใน details)
	count, _ := ae.alertRepo.CountRecentEvents(
		ctx,
		logEntry.TenantID,
		rule.EventType,
		rule.ConditionField,
		conditionValue,
		rule.TimeWindowSeconds,
	)

	// สร้าง alert details
	details := models.AlertDetails{
		"event_type":      rule.EventType,
		"condition_field": rule.ConditionField,
		"condition_value": conditionValue,
		"count":           count,
		"threshold":       rule.Threshold,
		"time_window":     rule.TimeWindowSeconds,
		"message":         fmt.Sprintf("%s detected: %s=%s occurred %d times (threshold: %d)", rule.RuleName, rule.ConditionField, conditionValue, count, rule.Threshold),
	}

	// เพิ่มข้อมูลเพิ่มเติมถ้ามี
	if srcIP, ok := bodyMap["src_ip"]; ok {
		details["src_ip"] = srcIP
	}
	if user, ok := bodyMap["user"]; ok {
		details["user"] = user
	}
	if host, ok := bodyMap["host"]; ok {
		details["host"] = host
	}

	// สร้าง alert
	alert := models.AlertHistory{
		TenantID:     logEntry.TenantID,
		RuleID:       rule.ID,
		RuleName:     rule.RuleName,
		TriggeredAt:  time.Now(),
		Details:      details,
		Severity:     rule.Severity,
		Acknowledged: false,
	}

	// บันทึกลง database
	if err := ae.alertRepo.CreateAlert(ctx, alert); err != nil {
		log.Printf("❌ Failed to create alert: %v", err)
		return
	}

	log.Printf("🚨 ALERT TRIGGERED: %s | %s=%s | count=%d | severity=%d",
		rule.RuleName, rule.ConditionField, conditionValue, count, rule.Severity)
}
