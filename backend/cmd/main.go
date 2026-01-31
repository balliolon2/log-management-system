package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq" // Postgres Driver

	"log-management-backend/internal/alert"
	"log-management-backend/internal/auth"
	"log-management-backend/internal/handler"
	"log-management-backend/internal/ingest"
	"log-management-backend/internal/job"
	"log-management-backend/internal/models"
	"log-management-backend/internal/normalize"
	"log-management-backend/internal/repository"
)

func main() {
	// ... (Database & Repo init omitted for brevity, assuming standard setup) ...

	connStr := "postgres://user:password@localhost:5432/logdb?sslmode=disable"
	if os.Getenv("DB_HOST") != "" {
		connStr = fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Init Repositories
	logRepo := repository.NewLogRepository(db)
	alertRepo := repository.NewAlertRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 3. Init Handlers
	logHandler := handler.NewLogHandler(logRepo)
	alertHandler := handler.NewAlertHandler(alertRepo)
	authHandler := handler.NewAuthHandler(userRepo)

	// 4. Init & Start Alert Engine
	alertEngine := alert.NewAlertEngine(alertRepo, logRepo)
	if err := alertEngine.Start(context.Background()); err != nil {
		log.Fatalf("Failed to start alert engine: %v", err)
	}

	// 5. Background Jobs (Syslog + Retention)
	go ingest.StartSyslogServerWithAlerts(context.Background(), logRepo, alertEngine, "514")
	go ingest.StartSyslogTCPServerWithAlerts(context.Background(), logRepo, alertEngine, "514")

	// Start Retention Job (7 days)
	go job.StartRetentionJob(context.Background(), logRepo, 7)

	// 6. CORS Preflight Handler (สำหรับ OPTIONS requests)
	corsHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	// 7. Setup Routes (API Endpoints)

	// Auth Routes (ไม่ต้องใช้ middleware)
	http.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		corsHandler(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		authHandler.LoginHandler(w, r)
	})

	http.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		corsHandler(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		authHandler.RegisterHandler(w, r)
	})

	// Protected Routes (ต้องใช้ middleware)
	http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		corsHandler(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		auth.AuthMiddleware(logHandler.GetLogsHandler)(w, r)
	})

	http.HandleFunc("/api/alerts", func(w http.ResponseWriter, r *http.Request) {
		corsHandler(w, r)
		if r.Method == http.MethodOptions {
			return
		}
		auth.AuthMiddleware(alertHandler.GetAlertsHandler)(w, r)
	})

	// Ingest Endpoint (ไม่ต้องใช้ auth - เพื่อให้ระบบภายนอกส่ง log ได้)
	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		corsHandler(w, r)
		if r.Method == http.MethodOptions {
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 1. อ่าน Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// 2. Normalize
		entries, err := normalize.ProcessJSONLog(body, "default_tenant")
		if err != nil {
			log.Printf("Normalize error: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// 3. Save to DB (Batch)
		if err := logRepo.CreateBatch(r.Context(), entries); err != nil {
			log.Printf("Failed to save HTTP log: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// 4. Check Alert (Real-time) - Loop for each entry
		go func(logs []models.LogEntry) {
			for _, entry := range logs {
				alertEngine.CheckLog(context.Background(), entry)
			}
		}(entries)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Received %d logs", len(entries))))
	})

	log.Println("🚀 HTTP Server starting on :8080")
	log.Println("🔐 Auth APIs:")
	log.Println("   POST http://localhost:8080/login")
	log.Println("   POST http://localhost:8080/register")
	log.Println("📊 Protected APIs (require JWT token):")
	log.Println("   GET http://localhost:8080/api/logs")
	log.Println("   GET http://localhost:8080/api/alerts")
	log.Println("📥 Ingest API:")
	log.Println("   POST http://localhost:8080/ingest")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
