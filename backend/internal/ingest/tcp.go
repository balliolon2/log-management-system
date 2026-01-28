package ingest

import (
	"bufio"
	"context"
	"log"
	"net"

	"log-management-backend/internal/alert"
	"log-management-backend/internal/models"
	"log-management-backend/internal/normalize"
	"log-management-backend/internal/repository"
)

// StartSyslogTCPServer - ฟังก์ชันเดิม (ไม่มี Alert)
func StartSyslogTCPServer(ctx context.Context, repo *repository.LogRepository, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start TCP syslog server: %v", err)
	}
	defer listener.Close()
	log.Printf("Syslog TCP server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCP accept error: %v", err)
			continue
		}
		go handleTCPConnection(ctx, conn, repo, nil) // nil = ไม่มี Alert Engine
	}
}

// StartSyslogTCPServerWithAlerts - ฟังก์ชันใหม่ (รองรับ Alert Engine)
func StartSyslogTCPServerWithAlerts(ctx context.Context, repo *repository.LogRepository, alertEngine *alert.AlertEngine, port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to start TCP syslog server: %v", err)
	}
	defer listener.Close()
	log.Printf("🔥 Syslog TCP server (with Alerts) listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCP accept error: %v", err)
			continue
		}
		go handleTCPConnection(ctx, conn, repo, alertEngine)
	}
}

// handleTCPConnection - แก้ไขให้รองรับ Alert Engine
func handleTCPConnection(ctx context.Context, conn net.Conn, repo *repository.LogRepository, alertEngine *alert.AlertEngine) {
	defer conn.Close()
	log.Printf("DEBUG: New TCP connection from %s", conn.RemoteAddr().String())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		raw := scanner.Text()
		log.Printf("DEBUG: TCP Scanner read: %s", raw)
		if raw == "" {
			continue
		}

		entry := normalize.Syslog(raw, "syslog_default")
		entry.Source = "firewall-tcp"

		if err := repo.CreateBatch(ctx, []models.LogEntry{entry}); err != nil {
			log.Printf("Failed to save TCP log: %v", err)
		} else {
			log.Println("DEBUG: TCP Log saved to DB")

			// ถ้ามี Alert Engine ให้ตรวจสอบ log
			if alertEngine != nil {
				go alertEngine.CheckLog(context.Background(), entry)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("TCP connection read error: %v", err)
	}
}
