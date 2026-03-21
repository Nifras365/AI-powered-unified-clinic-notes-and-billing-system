package repositories

import (
	"context"
	"fmt"

	"clinicnotes/backend/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConsultationDetails struct {
	Patient       models.Patient
	Consultation  models.Consultation
	Prescriptions []models.Prescription
	LabTests      []models.LabTest
	Billing       models.Billing
}

type ConsultationRepository struct {
	db *pgxpool.Pool
}

func NewConsultationRepository(db *pgxpool.Pool) *ConsultationRepository {
	return &ConsultationRepository{db: db}
}

func (r *ConsultationRepository) SaveConsultationBundle(ctx context.Context, details ConsultationDetails) (ConsultationDetails, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return ConsultationDetails{}, err
	}
	defer tx.Rollback(ctx)

	if details.Patient.ID == "" {
		details.Patient.ID = uuid.NewString()
	}
	if details.Consultation.ID == "" {
		details.Consultation.ID = uuid.NewString()
	}
	details.Consultation.PatientID = details.Patient.ID

	err = tx.QueryRow(ctx, `INSERT INTO patients(id, name, age) VALUES ($1, $2, $3) RETURNING created_at`, details.Patient.ID, details.Patient.Name, details.Patient.Age).Scan(&details.Patient.CreatedAt)
	if err != nil {
		return ConsultationDetails{}, err
	}

	err = tx.QueryRow(ctx, `INSERT INTO consultations(id, patient_id, raw_input, observations) VALUES ($1, $2, $3, $4) RETURNING created_at`, details.Consultation.ID, details.Patient.ID, details.Consultation.RawInput, details.Consultation.Observations).Scan(&details.Consultation.CreatedAt)
	if err != nil {
		return ConsultationDetails{}, err
	}

	for i, p := range details.Prescriptions {
		if p.ID == "" {
			p.ID = uuid.NewString()
		}
		p.ConsultationID = details.Consultation.ID
		details.Prescriptions[i] = p
		_, err := tx.Exec(ctx, `INSERT INTO prescriptions(id, consultation_id, drug_name, dosage, frequency, price) VALUES ($1, $2, $3, $4, $5, $6)`, p.ID, p.ConsultationID, p.DrugName, p.Dosage, p.Frequency, p.Price)
		if err != nil {
			return ConsultationDetails{}, err
		}
	}

	for i, t := range details.LabTests {
		if t.ID == "" {
			t.ID = uuid.NewString()
		}
		t.ConsultationID = details.Consultation.ID
		details.LabTests[i] = t
		_, err := tx.Exec(ctx, `INSERT INTO lab_tests(id, consultation_id, test_name, price) VALUES ($1, $2, $3, $4)`, t.ID, t.ConsultationID, t.TestName, t.Price)
		if err != nil {
			return ConsultationDetails{}, err
		}
	}

	if details.Billing.ID == "" {
		details.Billing.ID = uuid.NewString()
	}
	details.Billing.ConsultationID = details.Consultation.ID

	err = tx.QueryRow(ctx, `
		INSERT INTO billing(
			id, consultation_id, consultation_fee, drugs_total, tests_total, subtotal,
			discount_type, discount_value, discount_amount, total_amount
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING created_at
	`, details.Billing.ID, details.Billing.ConsultationID, details.Billing.ConsultationFee, details.Billing.DrugsTotal, details.Billing.TestsTotal, details.Billing.Subtotal, details.Billing.DiscountType, details.Billing.DiscountValue, details.Billing.DiscountAmount, details.Billing.TotalAmount).Scan(&details.Billing.CreatedAt)
	if err != nil {
		return ConsultationDetails{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return ConsultationDetails{}, err
	}

	return details, nil
}

func (r *ConsultationRepository) GetConsultationByID(ctx context.Context, consultationID string) (ConsultationDetails, error) {
	var out ConsultationDetails

	err := r.db.QueryRow(ctx, `
		SELECT p.id, p.name, p.age, p.created_at,
		       c.id, c.patient_id, c.raw_input, c.observations, c.created_at
		FROM consultations c
		JOIN patients p ON p.id = c.patient_id
		WHERE c.id = $1
	`, consultationID).Scan(
		&out.Patient.ID,
		&out.Patient.Name,
		&out.Patient.Age,
		&out.Patient.CreatedAt,
		&out.Consultation.ID,
		&out.Consultation.PatientID,
		&out.Consultation.RawInput,
		&out.Consultation.Observations,
		&out.Consultation.CreatedAt,
	)
	if err != nil {
		return ConsultationDetails{}, err
	}

	presRows, err := r.db.Query(ctx, `SELECT id, consultation_id, drug_name, dosage, frequency, price FROM prescriptions WHERE consultation_id = $1`, consultationID)
	if err != nil {
		return ConsultationDetails{}, err
	}
	defer presRows.Close()
	for presRows.Next() {
		var p models.Prescription
		if err := presRows.Scan(&p.ID, &p.ConsultationID, &p.DrugName, &p.Dosage, &p.Frequency, &p.Price); err != nil {
			return ConsultationDetails{}, err
		}
		out.Prescriptions = append(out.Prescriptions, p)
	}

	testRows, err := r.db.Query(ctx, `SELECT id, consultation_id, test_name, price FROM lab_tests WHERE consultation_id = $1`, consultationID)
	if err != nil {
		return ConsultationDetails{}, err
	}
	defer testRows.Close()
	for testRows.Next() {
		var t models.LabTest
		if err := testRows.Scan(&t.ID, &t.ConsultationID, &t.TestName, &t.Price); err != nil {
			return ConsultationDetails{}, err
		}
		out.LabTests = append(out.LabTests, t)
	}

	err = r.db.QueryRow(ctx, `
		SELECT id, consultation_id, consultation_fee, drugs_total, tests_total, subtotal,
		       discount_type, discount_value, discount_amount, total_amount, created_at
		FROM billing WHERE consultation_id = $1
	`, consultationID).Scan(
		&out.Billing.ID,
		&out.Billing.ConsultationID,
		&out.Billing.ConsultationFee,
		&out.Billing.DrugsTotal,
		&out.Billing.TestsTotal,
		&out.Billing.Subtotal,
		&out.Billing.DiscountType,
		&out.Billing.DiscountValue,
		&out.Billing.DiscountAmount,
		&out.Billing.TotalAmount,
		&out.Billing.CreatedAt,
	)
	if err != nil {
		return ConsultationDetails{}, err
	}

	return out, nil
}

func (r *ConsultationRepository) GetBillingByConsultationID(ctx context.Context, consultationID string) (models.Billing, error) {
	var b models.Billing
	err := r.db.QueryRow(ctx, `
		SELECT id, consultation_id, consultation_fee, drugs_total, tests_total, subtotal,
		       discount_type, discount_value, discount_amount, total_amount, created_at
		FROM billing WHERE consultation_id = $1
	`, consultationID).Scan(
		&b.ID,
		&b.ConsultationID,
		&b.ConsultationFee,
		&b.DrugsTotal,
		&b.TestsTotal,
		&b.Subtotal,
		&b.DiscountType,
		&b.DiscountValue,
		&b.DiscountAmount,
		&b.TotalAmount,
		&b.CreatedAt,
	)
	if err != nil {
		return models.Billing{}, fmt.Errorf("billing not found: %w", err)
	}
	return b, nil
}
