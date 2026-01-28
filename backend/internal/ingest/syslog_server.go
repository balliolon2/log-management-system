package ingest

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"log-management-backend/internal/alert"
	"log-management-backend/internal/models"
	"log-management-backend/internal/normalize"
	"log-management-backend/internal/repository"
)

// Config สำหรับ Batching
const (
	BatchSize    = 100             // เก็บครบ 100 logs แล้วค่อยบันทึก
	BatchTimeout = 2 * time.Second // หรือผ่านไป 2 วินาทีแล้วค่อยบันทึก
	BufferSize   = 1000            // ขนาดสายพาน (Channel)
)

// StartSyslogServer - ฟังก์ชันเดิม (ไม่มี Alert)
func StartSyslogServer(ctx context.Context, repo *repository.LogRepository, port string) {
	// 1. Setup UDP Listener (เหมือนเดิม)
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()
	log.Printf("Syslog UDP server listening on port %s", port)

	// 2. สร้างสายพาน (Channel) สำหรับส่ง Log ไปให้ Worker
	logChan := make(chan models.LogEntry, BufferSize)

	// 3. เริ่ม Worker (ทำงานเบื้องหลัง)
	go batchWorker(ctx, repo, logChan, nil) // nil = ไม่มี Alert Engine

	// 4. Loop รับของ (Producer)
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Error reading UDP: %v", err)
				continue
			}
			rawMsg := string(buf[:n])
			rawMsg = strings.TrimSpace(rawMsg)
			if rawMsg == "" {
				continue
			}

			// Debug
			log.Printf("DEBUG: Received UDP Log: %s", rawMsg)

			// Normalize
			entry := normalize.Syslog(rawMsg, "syslog_default")

			// ส่งเข้าสายพาน (Non-blocking ถ้า buffer ไม่เต็ม)
			select {
			case logChan <- entry:
				// ส่งสำเร็จ
			default:
				log.Printf("Log Buffer Full! Dropping log.")
			}
		}
	}
}

// StartSyslogServerWithAlerts - ฟังก์ชันใหม่ (รองรับ Alert Engine)
func StartSyslogServerWithAlerts(ctx context.Context, repo *repository.LogRepository, alertEngine *alert.AlertEngine, port string) {
	// 1. Setup UDP Listener
	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to resolve UDP address: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on UDP: %v", err)
	}
	defer conn.Close()
	log.Printf("🔥 Syslog UDP server (with Alerts) listening on port %s", port)

	// 2. สร้างสายพาน
	logChan := make(chan models.LogEntry, BufferSize)

	// 3. เริ่ม Worker พร้อม Alert Engine
	go batchWorker(ctx, repo, logChan, alertEngine)

	// 4. Loop รับของ
	buf := make([]byte, 4096)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Error reading UDP: %v", err)
				continue
			}
			rawMsg := string(buf[:n])
			rawMsg = strings.TrimSpace(rawMsg)
			if rawMsg == "" {
				continue
			}

			log.Printf("DEBUG: Received UDP Log: %s", rawMsg)
			entry := normalize.Syslog(rawMsg, "syslog_default")

			select {
			case logChan <- entry:
			default:
				log.Printf("Log Buffer Full! Dropping log.")
			}
		}
	}
}

// batchWorker - แก้ไขให้รองรับ Alert Engine
func batchWorker(ctx context.Context, repo *repository.LogRepository, logChan <-chan models.LogEntry, alertEngine *alert.AlertEngine) {
	batch := make([]models.LogEntry, 0, BatchSize)
	ticker := time.NewTicker(BatchTimeout)
	defer ticker.Stop()

	flush := func() {
		if len(batch) > 0 {
			if err := repo.CreateBatch(context.Background(), batch); err != nil {
				log.Printf("❌ Failed to insert batch: %v", err)
			} else {
				log.Printf("✅ Batch saved: %d logs", len(batch))

				// ถ้ามี Alert Engine ให้ตรวจสอบแต่ละ log
				if alertEngine != nil {
					for _, entry := range batch {
						go alertEngine.CheckLog(context.Background(), entry)
					}
				}
			}
			batch = batch[:0]
		}
	}

	for {
		select {
		case entry := <-logChan:
			batch = append(batch, entry)
			if len(batch) >= BatchSize {
				flush()
			}
		case <-ticker.C:
			flush()
		case <-ctx.Done():
			flush()
			return
		}
	}
}
