package routes

import (
	"clinicnotes/backend/internal/handlers"
	"clinicnotes/backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Health       *handlers.HealthHandler
	Parse        *handlers.ParseHandler
	Consultation *handlers.ConsultationHandler
	Billing      *handlers.BillingHandler
	Report       *handlers.ReportHandler
}

func Register(router *gin.Engine, h Handlers) {
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(middleware.ErrorHandler())

	router.GET("/health", h.Health.Health)

	api := router.Group("/api")
	{
		api.POST("/parse", h.Parse.Parse)
		api.POST("/consultations", h.Consultation.Create)
		api.GET("/consultations/:id", h.Consultation.GetByID)
		api.GET("/billing/:id", h.Billing.GetByConsultationID)
		api.GET("/reports/:id", h.Report.GetByConsultationID)
	}
}
