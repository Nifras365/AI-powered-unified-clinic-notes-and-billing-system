package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"clinicnotes/backend/internal/config"
	"clinicnotes/backend/internal/handlers"
	"clinicnotes/backend/internal/repositories"
	"clinicnotes/backend/internal/routes"
	"clinicnotes/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	if err := applyMigrations(ctx, db); err != nil {
		log.Fatal(err)
	}

	gemini := services.NewGeminiService(cfg.GeminiAPIKey, cfg.GeminiModel)
	ai := services.NewAIOrchestrator(gemini, nil)

	repo := repositories.NewConsultationRepository(db)
	billingService := services.NewBillingService(cfg.BaseConsultationFee, cfg.DefaultDiscountType, cfg.DefaultDiscountValue)
	reportService := services.NewReportService()

	healthHandler := handlers.NewHealthHandler()
	parseHandler := handlers.NewParseHandler(ai)
	consultationHandler := handlers.NewConsultationHandler(repo, billingService)
	billingHandler := handlers.NewBillingHandler(repo)
	reportHandler := handlers.NewReportHandler(repo, reportService)

	r := gin.Default()
	routes.Register(r, routes.Handlers{
		Health:       healthHandler,
		Parse:        parseHandler,
		Consultation: consultationHandler,
		Billing:      billingHandler,
		Report:       reportHandler,
	})

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func applyMigrations(ctx context.Context, db *pgxpool.Pool) error {
	migrationPath := filepath.Join("migrations", "001_init.sql")
	content, err := os.ReadFile(migrationPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, string(content))
	return err
}
