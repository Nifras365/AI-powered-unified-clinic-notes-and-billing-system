package dto

type Drug struct {
	Name      string `json:"name"`
	Dosage    string `json:"dosage"`
	Frequency string `json:"frequency"`
}

type ParsedResult struct {
	Drugs        []Drug   `json:"drugs"`
	LabTests     []string `json:"lab_tests"`
	Observations string   `json:"observations"`
}

type ParseRequest struct {
	RawInput string `json:"raw_input"`
}

type ParseResponse struct {
	Parsed ParsedResult `json:"parsed"`
}

type PatientInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type PricingInput struct {
	ConsultationFee float64            `json:"consultation_fee"`
	DiscountType    string             `json:"discount_type"`
	DiscountValue   float64            `json:"discount_value"`
	DrugPrices      map[string]float64 `json:"drug_prices"`
	TestPrices      map[string]float64 `json:"test_prices"`
}

type CreateConsultationRequest struct {
	Patient  PatientInput `json:"patient"`
	RawInput string       `json:"raw_input"`
	Parsed   ParsedResult `json:"parsed"`
	Pricing  PricingInput `json:"pricing"`
}

type ConsultationResponse struct {
	ID        string       `json:"id"`
	PatientID string       `json:"patient_id"`
	RawInput  string       `json:"raw_input"`
	Parsed    ParsedResult `json:"parsed"`
	CreatedAt string       `json:"created_at"`
}
