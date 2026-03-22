# AI-Powered Unified Clinic Notes and Billing System

Production-ready monorepo for capturing clinical notes, extracting structured entities with AI, enforcing human verification, calculating billing, and generating printable/downloadable reports.

## 1. Project Lifecycle Documentation

### 1.1 Problem and Goals

Clinicians need a single workflow to:
- Capture clinical notes quickly (typed or voice)
- Convert unstructured notes into structured medical data
- Verify AI output before persistence
- Generate billing and final reports in one flow

Primary goals:
- Verification-first medical data handling
- Transaction-safe persistence across related entities
- Simple and fast user workflow
- Printable and downloadable report output

### 1.2 Scope and Boundaries

In scope:
- AI parsing of medical notes into strict JSON
- Manual review/edit before saving
- Billing with consultation fee, medicine/test prices, discount
- Report preview, print, and PDF download

Out of scope (current version):
- Authentication and role-based access control
- Master pricing catalog management
- Multi-clinic tenant support

### 1.3 Lifecycle Stages Implemented

1. Requirement definition
2. System design and architecture
3. Database schema design and migration
4. Backend implementation (API + business logic)
5. Frontend implementation (verification-first UI)
6. End-to-end integration
7. Build validation and final documentation

## 2. System Design

### 2.1 Functional Flow

1. Doctor enters notes (text or voice)
2. Frontend sends notes to parse endpoint
3. Backend calls AI provider and validates strict schema
4. Frontend displays editable parsed entities with highlights
5. User verifies/edits parsed data
6. Confirm save triggers transactional persistence and billing
7. Billing + report endpoints return finalized outputs
8. Report can be previewed, printed, or downloaded as PDF

### 2.2 Non-Functional Design

- Reliability: transactional write model for consistency
- Data quality: mandatory human verification before persistence
- Maintainability: layered architecture (handlers/services/repository)
- UX speed: responsive frontend with single-screen workflow

## 4. Implementation Details

### 4.1 Backend (Go + Gin)

Key modules:
- Config loading and defaults
- Parse, consultation, billing, and report handlers
- AI orchestration and Gemini integration
- Billing service and report rendering service
- Repository with transaction-safe save operations

API endpoints:
- `GET /health`
- `POST /api/parse`
- `POST /api/consultations`
- `GET /api/consultations/:id`
- `GET /api/billing/:id`
- `GET /api/reports/:id`

Important behavior:
- Parse endpoint enforces strict JSON structure validation
- Save endpoint persists patient + consultation + prescriptions + tests + 

### 4.2 Frontend (React + TypeScript + Vite)

UI components:
- InputPanel: text/voice capture
- HighlightedTextViewer: visual entity highlighting
- ParsedResultsEditor: mandatory verification and edits
- BillingView: finalized billing summary
- ReportView: preview + print + PDF download

Current billing controls:
- Consultation Fee
- Discount Type (`fixed` or `percentage`)
- Discount Value

Report actions:
- Print
- Download PDF

### 4.3 AI Prompting and Validation Strategy

- Fixed output schema contract
- Markdown/prose stripping from model output
- JSON extraction and decode
- Structural validation with retries
- Reject invalid output after retry budget

## 5. Database Scripts (Tables + Relationships)

Primary migration script (`backend/migrations/001_init.sql`):

```sql
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
```

Relationship summary:
- One `patients` row can have many `consultations`
- One `consultations` row can have many `prescriptions`
- One `consultations` row can have many `lab_tests`
- One `consultations` row has exactly one `billing` row
- Cascading deletes are enabled from consultation-parent relationships

## 6. API Contract Snapshot

### 6.1 Parse

`POST /api/parse`

Request:

```json
{
  "raw_input": "Patient has fever. Start paracetamol 500mg twice daily and order CBC."
}
```

Response:

```json
{
  "parsed": {
    "drugs": [
      { "name": "Paracetamol", "dosage": "500mg", "frequency": "twice daily" }
    ],
    "lab_tests": ["CBC"],
    "observations": "Patient has fever"
  }
}
```

### 6.2 Save Consultation

`POST /api/consultations`

Request:

```json
{
  "patient": { "name": "John Doe", "age": 34 },
  "raw_input": "...",
  "parsed": {
    "drugs": [{ "name": "Paracetamol", "dosage": "500mg", "frequency": "BID" }],
    "lab_tests": ["CBC"],
    "observations": "Fever"
  },
  "pricing": {
    "consultation_fee": 40,
    "discount_type": "fixed",
    "discount_value": 2,
    "drug_prices": { "Paracetamol": 5 },
    "test_prices": { "CBC": 10 }
  }
}
```

### 6.3 Fetch Billing and Report

- `GET /api/billing/:id`
- `GET /api/reports/:id` (returns `text/html`)

## 7. Setup and Run

### 7.1 Prerequisites

- Go 1.23+
- Node.js 20+
- PostgreSQL 16+

### 7.2 Backend

```bash
cd backend
cp .env.example .env
go mod tidy
go run ./cmd/server
```

### 7.3 Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Defaults:
- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8080`

## 8. Final Outcomes

Delivered outcomes:
- End-to-end AI-assisted consultation workflow
- Mandatory verification before persistence
- Transaction-safe multi-entity save
- Billing generation with discount handling
- Printable report preview
- Report export options: Print and Download PDF

Technical validation outcomes:
- Frontend build compiles successfully
- Backend module build compiles successfully

## 9. Environment Variables

Backend:
- `APP_ENV`
- `PORT`
- `DATABASE_URL`
- `GEMINI_API_KEY`
- `GEMINI_MODEL`
- `BASE_CONSULTATION_FEE`
- `DEFAULT_DISCOUNT_TYPE`
- `DEFAULT_DISCOUNT_VALUE`

Frontend:
- `VITE_API_BASE_URL`
