package handlers

import (
	"net/http"
	"strings"

	"clinicnotes/backend/internal/dto"
	"clinicnotes/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ParseHandler struct {
	ai *services.AIOrchestrator
}

func NewParseHandler(ai *services.AIOrchestrator) *ParseHandler {
	return &ParseHandler{ai: ai}
}

func (h *ParseHandler) Parse(c *gin.Context) {
	var req dto.ParseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	req.RawInput = strings.TrimSpace(req.RawInput)
	if req.RawInput == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "raw_input is required"})
		return
	}

	parsed, err := h.ai.ParseMedicalText(c.Request.Context(), req.RawInput)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.ParseResponse{Parsed: parsed})
}
