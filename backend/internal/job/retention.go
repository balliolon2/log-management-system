package job

import (
	"context"
	"log"
	"time"

	"log-management-backend/internal/repository"
)

// StartRetentionJob starts a background ticker to delete old logs
func StartRetentionJob(ctx context.Context, repo *repository.LogRepository, retentionDays int) {
	log.Printf("🕰️ Retention Job started: Delete logs older than %d days", retentionDays)

	// Run immediately on start
	runCleanup(ctx, repo, retentionDays)

	ticker := time.NewTicker(1 * time.Hour) // Check every hour
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Retention Job stopped")
			return
		case <-ticker.C:
			runCleanup(ctx, repo, retentionDays)
		}
	}
}

func runCleanup(ctx context.Context, repo *repository.LogRepository, days int) {
	count, err := repo.DeleteOldLogs(ctx, days)
	if err != nil {
		log.Printf("❌ Retention Job failed: %v", err)
	} else if count > 0 {
		log.Printf("🧹 Retention Job: Deleted %d old logs", count)
	}
}
