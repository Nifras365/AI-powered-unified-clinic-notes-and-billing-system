package models

import "time"

type Patient struct {
	ID        string
	Name      string
	Age       int
	CreatedAt time.Time
}

type Consultation struct {
	ID           string
	PatientID    string
	RawInput     string
	Observations string
	CreatedAt    time.Time
}

type Prescription struct {
	ID             string
	ConsultationID string
	DrugName       string
	Dosage         string
	Frequency      string
	Price          float64
}

type LabTest struct {
	ID             string
	ConsultationID string
	TestName       string
	Price          float64
}

type Billing struct {
	ID              string    `json:"id"`
	ConsultationID  string    `json:"consultation_id"`
	ConsultationFee float64   `json:"consultation_fee"`
	DrugsTotal      float64   `json:"drugs_total"`
	TestsTotal      float64   `json:"tests_total"`
	Subtotal        float64   `json:"subtotal"`
	DiscountType    string    `json:"discount_type"`
	DiscountValue   float64   `json:"discount_value"`
	DiscountAmount  float64   `json:"discount_amount"`
	TotalAmount     float64   `json:"total_amount"`
	CreatedAt       time.Time `json:"created_at"`
}
