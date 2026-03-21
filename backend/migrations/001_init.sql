CREATE TABLE IF NOT EXISTS patients (
    id UUID PRIMARY KEY,
    name VARCHAR(150) NOT NULL,
    age INT NOT NULL CHECK (age > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS consultations (
    id UUID PRIMARY KEY,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    raw_input TEXT NOT NULL,
    observations TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS prescriptions (
    id UUID PRIMARY KEY,
    consultation_id UUID NOT NULL REFERENCES consultations(id) ON DELETE CASCADE,
    drug_name VARCHAR(150) NOT NULL,
    dosage VARCHAR(100) NOT NULL,
    frequency VARCHAR(100) NOT NULL,
    price NUMERIC(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS lab_tests (
    id UUID PRIMARY KEY,
    consultation_id UUID NOT NULL REFERENCES consultations(id) ON DELETE CASCADE,
    test_name VARCHAR(150) NOT NULL,
    price NUMERIC(12,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS billing (
    id UUID PRIMARY KEY,
    consultation_id UUID NOT NULL UNIQUE REFERENCES consultations(id) ON DELETE CASCADE,
    consultation_fee NUMERIC(12,2) NOT NULL,
    drugs_total NUMERIC(12,2) NOT NULL,
    tests_total NUMERIC(12,2) NOT NULL,
    subtotal NUMERIC(12,2) NOT NULL,
    discount_type VARCHAR(20) NOT NULL,
    discount_value NUMERIC(12,2) NOT NULL,
    discount_amount NUMERIC(12,2) NOT NULL,
    total_amount NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_consultations_patient_id ON consultations(patient_id);
CREATE INDEX IF NOT EXISTS idx_prescriptions_consultation_id ON prescriptions(consultation_id);
CREATE INDEX IF NOT EXISTS idx_lab_tests_consultation_id ON lab_tests(consultation_id);
