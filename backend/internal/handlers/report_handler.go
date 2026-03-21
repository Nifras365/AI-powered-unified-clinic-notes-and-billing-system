package handlers

import (
	"net/http"

	"clinicnotes/backend/internal/repositories"
	"clinicnotes/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	repo          *repositories.ConsultationRepository
	reportService *services.ReportService
}

func NewReportHandler(repo *repositories.ConsultationRepository, reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{repo: repo, reportService: reportService}
}

func (h *ReportHandler) GetByConsultationID(c *gin.Context) {
	id := c.Param("id")
	details, err := h.repo.GetConsultationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "report data not found"})
		return
	}

	html, err := h.reportService.BuildHTMLReport(services.ReportPayload{
		Patient:       details.Patient,
		Consultation:  details.Consultation,
		Prescriptions: details.Prescriptions,
		LabTests:      details.LabTests,
		Billing:       details.Billing,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
