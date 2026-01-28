package normalize

import (
	"encoding/json"
	"testing"
)

func TestProcessJSONLog(t *testing.T) {
	// Case 1: Single Object
	singleJSON := `{"event_type": "login", "user": "admin"}`
	entries, err := ProcessJSONLog([]byte(singleJSON), "test_tenant")
	if err != nil {
		t.Fatalf("Failed to process single JSON: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(entries))
	}
	if entries[0].EventType != "login" {
		t.Errorf("Expected event_type 'login', got '%s'", entries[0].EventType)
	}

	// Case 2: Batch Array
	batchJSON := `[
		{"event_type": "login", "user": "u1"},
		{"event_type": "logout", "user": "u1"}
	]`
	entries, err = ProcessJSONLog([]byte(batchJSON), "test_tenant")
	if err != nil {
		t.Fatalf("Failed to process batch JSON: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
	if entries[0].EventType != "login" {
		t.Errorf("Expected first event 'login', got '%s'", entries[0].EventType)
	}
	if entries[1].EventType != "logout" {
		t.Errorf("Expected second event 'logout', got '%s'", entries[1].EventType)
	}

	// Verify Body content
	var bodyMap map[string]interface{}
	json.Unmarshal(entries[0].Body, &bodyMap)
	if bodyMap["user"] != "u1" {
		t.Errorf("Expected user 'u1', got '%v'", bodyMap["user"])
	}
}
