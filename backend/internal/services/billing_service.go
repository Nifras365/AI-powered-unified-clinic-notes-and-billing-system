package services

import (
	"fmt"
	"strings"

	"clinicnotes/backend/internal/dto"
	"clinicnotes/backend/internal/models"
)

type BillingService struct {
	defaultConsultationFee float64
	defaultDiscountType    string
	defaultDiscountValue   float64
}

func NewBillingService(defaultConsultationFee float64, defaultDiscountType string, defaultDiscountValue float64) *BillingService {
	return &BillingService{
		defaultConsultationFee: defaultConsultationFee,
		defaultDiscountType:    defaultDiscountType,
		defaultDiscountValue:   defaultDiscountValue,
	}
}

func (s *BillingService) BuildBilling(parsed dto.ParsedResult, pricing dto.PricingInput, consultationID string) (models.Billing, []models.Prescription, []models.LabTest, error) {
	consultationFee := pricing.ConsultationFee
	if consultationFee <= 0 {
		consultationFee = s.defaultConsultationFee
	}

	discountType := strings.ToLower(strings.TrimSpace(pricing.DiscountType))
	if discountType == "" {
		discountType = strings.ToLower(s.defaultDiscountType)
	}
	discountValue := pricing.DiscountValue
	if discountValue < 0 {
		discountValue = s.defaultDiscountValue
	}

	prescriptions := make([]models.Prescription, 0, len(parsed.Drugs))
	labTests := make([]models.LabTest, 0, len(parsed.LabTests))

	drugsTotal := 0.0
	for _, d := range parsed.Drugs {
		price := resolvePrice(pricing.DrugPrices, d.Name)
		drugsTotal += price
		prescriptions = append(prescriptions, models.Prescription{
			ConsultationID: consultationID,
			DrugName:       d.Name,
			Dosage:         d.Dosage,
			Frequency:      d.Frequency,
			Price:          price,
		})
	}

	testsTotal := 0.0
	for _, t := range parsed.LabTests {
		price := resolvePrice(pricing.TestPrices, t)
		testsTotal += price
		labTests = append(labTests, models.LabTest{
			ConsultationID: consultationID,
			TestName:       t,
			Price:          price,
		})
	}

	subtotal := consultationFee + drugsTotal + testsTotal
	discountAmount, err := calculateDiscount(discountType, discountValue, subtotal)
	if err != nil {
		return models.Billing{}, nil, nil, err
	}
	total := subtotal - discountAmount
	if total < 0 {
		total = 0
	}

	billing := models.Billing{
		ConsultationID:  consultationID,
		ConsultationFee: consultationFee,
		DrugsTotal:      drugsTotal,
		TestsTotal:      testsTotal,
		Subtotal:        subtotal,
		DiscountType:    discountType,
		DiscountValue:   discountValue,
		DiscountAmount:  discountAmount,
		TotalAmount:     total,
	}

	return billing, prescriptions, labTests, nil
}

func calculateDiscount(discountType string, discountValue, subtotal float64) (float64, error) {
	switch discountType {
	case "fixed":
		if discountValue > subtotal {
			return subtotal, nil
		}
		return discountValue, nil
	case "percentage":
		if discountValue > 100 {
			discountValue = 100
		}
		return subtotal * (discountValue / 100), nil
	default:
		return 0, fmt.Errorf("invalid discount type: %s", discountType)
	}
}

func resolvePrice(priceMap map[string]float64, key string) float64 {
	if priceMap == nil {
		return 0
	}
	for k, v := range priceMap {
		if strings.EqualFold(strings.TrimSpace(k), strings.TrimSpace(key)) {
			if v < 0 {
				return 0
			}
			return v
		}
	}
	return 0
}
