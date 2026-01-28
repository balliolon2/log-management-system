package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"log-management-backend/internal/auth"
	"log-management-backend/internal/repository"
)

type LogHandler struct {
	repo *repository.LogRepository
}

func NewLogHandler(repo *repository.LogRepository) *LogHandler {
	return &LogHandler{repo: repo}
}

// GetLogsHandler - API สำหรับดึง logs (รองรับ tenant filter)
func (h *LogHandler) GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ดึง claims จาก context (มาจาก middleware)
	claims, ok := auth.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// ดึง query parameters
	query := r.URL.Query().Get("q")
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
		// ถ้ามี query parameter tenant ให้ใช้ตามที่ระบุ
		if t := r.URL.Query().Get("tenant"); t != "" {
			tenantID = t
		} else {
			tenantID = "" // ไม่ filter tenant (ดูทั้งหมด)
		}
	}
	// Viewer ดูได้เฉพาะ tenant ของตัวเอง (ใช้ tenantID จาก claims)

	// ดึงข้อมูล logs จาก database
	logs, err := h.repo.GetLogs(r.Context(), tenantID, query, limit)
	if err != nil {
		log.Printf("Failed to get logs: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// ถ้าไม่มีข้อมูล ส่ง empty array
	if logs == nil {
		logs = []map[string]interface{}{}
	}

	// ส่งกลับเป็น JSON
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(logs); err != nil {
		log.Printf("Failed to encode logs: %v", err)
	}
}
