package normalize

import (
	"encoding/json"
	"strings"
	"time"

	"log-management-backend/internal/models"
)

// Syslog แปลงข้อความ Syslog
func Syslog(raw string, tenantID string) models.LogEntry {
	// 1. Default Severity = 1 (Info)
	severity := 1

	// 2. ถ้าเจอคำว่า deny, block, down ให้เพิ่มระดับความรุนแรง
	lowerRaw := strings.ToLower(raw)
	if strings.Contains(lowerRaw, "deny") || strings.Contains(lowerRaw, "block") {
		severity = 5 // Warning
	} else if strings.Contains(lowerRaw, "link-down") || strings.Contains(lowerRaw, "error") {
		severity = 7 // Error
	}

	timestamp := time.Now()

	// แกะ body
	body := make(map[string]interface{})
	parts := strings.Fields(raw)
	for _, part := range parts {
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				body[kv[0]] = kv[1]
			}
		}
	}
	bodyBytes, _ := json.Marshal(body)

	return models.LogEntry{
		TenantID:   tenantID,
		Timestamp:  timestamp,
		Source:     "firewall",
		Severity:   severity,
		Body:       json.RawMessage(bodyBytes),
		RawMessage: raw,
		EventType:  "syslog_event",
	}
}

// ProcessJSONLog processes raw JSON input which can be a single object or an array of objects
func ProcessJSONLog(raw []byte, defaultTenant string) ([]models.LogEntry, error) {
	trimmed := strings.TrimSpace(string(raw))
	
	// Case 1: JSON Array (Batch)
	if strings.HasPrefix(trimmed, "[") {
		var rawList []map[string]interface{}
		if err := json.Unmarshal(raw, &rawList); err != nil {
			return nil, err
		}

		var entries []models.LogEntry
		for _, tempMap := range rawList {
			// Convert tempMap back to json bytes for RawMessage (approximation)
			// Or we can just keep the whole batch raw? 
			// Better: Let mapToLogEntry handle it using the map content
			entry := mapToLogEntry(tempMap, defaultTenant, "")
			entries = append(entries, entry)
		}
		return entries, nil
	}

	// Case 2: Single JSON Internal Object
	var tempMap map[string]interface{}
	if err := json.Unmarshal(raw, &tempMap); err != nil {
		return nil, err
	}
	
	entry := mapToLogEntry(tempMap, defaultTenant, string(raw))
	return []models.LogEntry{entry}, nil
}

// Deprecated: Use ProcessJSONLog instead
func JSON(raw []byte, defaultTenant string) (models.LogEntry, error) {
	entries, err := ProcessJSONLog(raw, defaultTenant)
	if err != nil {
		return models.LogEntry{}, err
	}
	if len(entries) == 0 {
		return models.LogEntry{}, nil
	}
	return entries[0], nil
}

// mapToLogEntry converts a raw map to LogEntry model
func mapToLogEntry(tempMap map[string]interface{}, defaultTenant string, rawString string) models.LogEntry {
	// 1. ดึง Common Fields
	ts := time.Now()
	if tStr, ok := tempMap["@timestamp"].(string); ok {
		if parsed, err := time.Parse(time.RFC3339, tStr); err == nil {
			ts = parsed
		}
	}

	tenant := defaultTenant
	if t, ok := tempMap["tenant"].(string); ok {
		tenant = t
	}

	source := "api"
	if s, ok := tempMap["source"].(string); ok {
		source = s
	}

	evtType := "unknown"
	if et, ok := tempMap["event_type"].(string); ok {
		evtType = et
	}

	severity := 1 // Default Info
	if sev, ok := tempMap["severity"].(float64); ok {
		severity = int(sev)
	}

	// Auto-Detect Severity
	if strings.Contains(strings.ToLower(evtType), "failed") ||
		strings.Contains(strings.ToLower(evtType), "malware") ||
		strings.Contains(strings.ToLower(evtType), "block") {
		if severity < 5 {
			severity = 7
		}
	}

	// *** สำคัญที่สุด: สร้าง body ที่ clean (ไม่มี metadata) ***
	body := make(map[string]interface{})

	// ถ้ามี field "body" ใน request → ใช้ body นั้น
	if bodyField, ok := tempMap["body"].(map[string]interface{}); ok {
		body = bodyField
	} else {
		// ถ้าไม่มี "body" → copy ทุก field ยกเว้น metadata
		for k, v := range tempMap {
			// ข้าม metadata fields
			if k != "tenant" && k != "@timestamp" && k != "timestamp" &&
				k != "source" && k != "event_type" && k != "severity" {
				body[k] = v
			}
		}
	}

	// เพิ่ม backward compatibility: ถ้ามี "ip" แต่ไม่มี "src_ip" → copy ไปให้
	if ip, ok := body["ip"].(string); ok {
		if _, hasSrcIP := body["src_ip"]; !hasSrcIP {
			body["src_ip"] = ip
		}
	}

	bodyBytes, _ := json.Marshal(body)
	
	// If rawString is empty (from batch), recreate it from tempMap
	if rawString == "" {
		rawBytes, _ := json.Marshal(tempMap)
		rawString = string(rawBytes)
	}

	return models.LogEntry{
		TenantID:   tenant,
		Timestamp:  ts,
		Source:     source,
		Severity:   severity,
		EventType:  evtType,
		Body:       json.RawMessage(bodyBytes),
		RawMessage: rawString,
	}
}
