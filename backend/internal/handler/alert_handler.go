package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"log-management-backend/internal/auth"
	"log-management-backend/internal/models"
	"log-management-backend/internal/repository"
)

type AlertHandler struct {
	alertRepo *repository.AlertRepository
}

func NewAlertHandler(alertRepo *repository.AlertRepository) *AlertHandler {
	return &AlertHandler{alertRepo: alertRepo}
}

// GetAlertsHandler - API สำหรับดึงประวัติ alerts (รองรับ tenant filter)
func (h *AlertHandler) GetAlertsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ดึง claims จาก context
	claims, ok := auth.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ดึง query parameters
	limitStr := r.URL.Query().Get("limit")
	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	// กำหนด tenant_id ตาม role
	tenantID := claims.TenantID
	if claims.Role == "admin" {
		// Admin ดูได้ทุก tenant
		if t := r.URL.Query().Get("tenant"); t != "" {
			tenantID = t
		} else {
			tenantID = "demo" // default สำหรับ admin
		}
	}
	// Viewer ดูได้เฉพาะ tenant ของตัวเอง

	// ดึงข้อมูล alerts จาก database
	alerts, err := h.alertRepo.GetAlertHistory(r.Context(), tenantID, limit)
	if err != nil {
		log.Printf("Failed to get alerts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ถ้าไม่มีข้อมูล ส่ง empty array
	if alerts == nil {
		alerts = []models.AlertHistory{}
	}

	// ส่งกลับเป็น JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(alerts); err != nil {
		log.Printf("Failed to encode alerts: %v", err)
	}
}
