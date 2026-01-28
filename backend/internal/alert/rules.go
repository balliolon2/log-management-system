package alert

import (
	"context"
	"log"

	"log-management-backend/internal/models"
	"log-management-backend/internal/repository"
)

type RuleManager struct {
	alertRepo *repository.AlertRepository
	rules     []models.AlertRule
}

func NewRuleManager(alertRepo *repository.AlertRepository) *RuleManager {
	return &RuleManager{
		alertRepo: alertRepo,
		rules:     []models.AlertRule{},
	}
}

// LoadRules - โหลดกฎจาก database
func (rm *RuleManager) LoadRules(ctx context.Context) error {
	rules, err := rm.alertRepo.GetAlertRules(ctx)
	if err != nil {
		return err
	}
	rm.rules = rules
	log.Printf("✅ Loaded %d alert rules", len(rm.rules))
	return nil
}

// GetRules - ดึงกฎทั้งหมด
func (rm *RuleManager) GetRules() []models.AlertRule {
	return rm.rules
}

// GetRuleByEventType - ดึงกฎที่ตรงกับ event_type
func (rm *RuleManager) GetRuleByEventType(eventType string) []models.AlertRule {
	var matched []models.AlertRule
	for _, rule := range rm.rules {
		if rule.EventType == eventType {
			matched = append(matched, rule)
		}
	}
	return matched
}
