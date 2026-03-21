package handlers

import (
	"net/http"

	"clinicnotes/backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type BillingHandler struct {
	repo *repositories.ConsultationRepository
}

func NewBillingHandler(repo *repositories.ConsultationRepository) *BillingHandler {
	return &BillingHandler{repo: repo}
}

func (h *BillingHandler) GetByConsultationID(c *gin.Context) {
	id := c.Param("id")
	billing, err := h.repo.GetBillingByConsultationID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "billing not found"})
		return
	}
	c.JSON(http.StatusOK, billing)
}
