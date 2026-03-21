package handlers

import (
	"net/http"
	"strings"

	"clinicnotes/backend/internal/dto"
	"clinicnotes/backend/internal/models"
	"clinicnotes/backend/internal/repositories"
	"clinicnotes/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type ConsultationHandler struct {
	repo           *repositories.ConsultationRepository
	billingService *services.BillingService
}

func NewConsultationHandler(repo *repositories.ConsultationRepository, billingService *services.BillingService) *ConsultationHandler {
	return &ConsultationHandler{repo: repo, billingService: billingService}
}

func (h *ConsultationHandler) Create(c *gin.Context) {
	var req dto.CreateConsultationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := validateCreateConsultation(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	billing, prescriptions, tests, err := h.billingService.BuildBilling(req.Parsed, req.Pricing, "")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	details := repositories.ConsultationDetails{
		Patient: models.Patient{
			Name: req.Patient.Name,
			Age:  req.Patient.Age,
		},
		Consultation: models.Consultation{
			RawInput:     req.RawInput,
			Observations: req.Parsed.Observations,
		},
		Prescriptions: prescriptions,
		LabTests:      tests,
		Billing:       billing,
	}

	saved, err := h.repo.SaveConsultationBundle(c.Request.Context(), details)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.ConsultationResponse{
		ID:        saved.Consultation.ID,
		PatientID: saved.Patient.ID,
		RawInput:  saved.Consultation.RawInput,
		Parsed: dto.ParsedResult{
			Drugs:        toDTODrugs(saved.Prescriptions),
			LabTests:     toDTOLabTests(saved.LabTests),
			Observations: saved.Consultation.Observations,
		},
		CreatedAt: saved.Consultation.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, response)
}

func (h *ConsultationHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	details, err := h.repo.GetConsultationByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "consultation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         details.Consultation.ID,
		"patient":    details.Patient,
		"raw_input":  details.Consultation.RawInput,
		"parsed":     gin.H{"drugs": toDTODrugs(details.Prescriptions), "lab_tests": toDTOLabTests(details.LabTests), "observations": details.Consultation.Observations},
		"created_at": details.Consultation.CreatedAt,
	})
}

func validateCreateConsultation(req dto.CreateConsultationRequest) error {
	if strings.TrimSpace(req.Patient.Name) == "" {
		return errInvalid("patient.name is required")
	}
	if req.Patient.Age <= 0 {
		return errInvalid("patient.age must be > 0")
	}
	if strings.TrimSpace(req.RawInput) == "" {
		return errInvalid("raw_input is required")
	}
	for _, d := range req.Parsed.Drugs {
		if strings.TrimSpace(d.Name) == "" {
			return errInvalid("drug name cannot be empty")
		}
	}
	return nil
}

func toDTODrugs(in []models.Prescription) []dto.Drug {
	out := make([]dto.Drug, 0, len(in))
	for _, p := range in {
		out = append(out, dto.Drug{
			Name:      p.DrugName,
			Dosage:    p.Dosage,
			Frequency: p.Frequency,
		})
	}
	return out
}

func toDTOLabTests(in []models.LabTest) []string {
	out := make([]string, 0, len(in))
	for _, t := range in {
		out = append(out, t.TestName)
	}
	return out
}

type invalidError struct{ msg string }

func (e invalidError) Error() string { return e.msg }

func errInvalid(msg string) error { return invalidError{msg: msg} }
