package services

import (
	"bytes"
	"html/template"

	"clinicnotes/backend/internal/models"
)

type ReportPayload struct {
	Patient       models.Patient
	Consultation  models.Consultation
	Prescriptions []models.Prescription
	LabTests      []models.LabTest
	Billing       models.Billing
}

type ReportService struct{}

func NewReportService() *ReportService {
	return &ReportService{}
}

func (s *ReportService) BuildHTMLReport(data ReportPayload) (string, error) {
	const tpl = `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8" />
<title>Clinic Report</title>
<style>
body { font-family: Arial, sans-serif; margin: 24px; color: #1f2937; }
h1,h2 { margin-bottom: 8px; }
table { width: 100%; border-collapse: collapse; margin-bottom: 16px; }
th,td { border: 1px solid #d1d5db; padding: 8px; text-align: left; }
.section { margin-bottom: 20px; }
</style>
</head>
<body>
  <h1>Consultation Report</h1>
  <div class="section">
    <h2>Patient Details</h2>
    <p><strong>Name:</strong> {{.Patient.Name}}</p>
    <p><strong>Age:</strong> {{.Patient.Age}}</p>
  </div>

  <div class="section">
    <h2>Prescription Report</h2>
    <table>
      <thead><tr><th>Drug</th><th>Dosage</th><th>Frequency</th><th>Price</th></tr></thead>
      <tbody>
      {{range .Prescriptions}}
        <tr><td>{{.DrugName}}</td><td>{{.Dosage}}</td><td>{{.Frequency}}</td><td>{{printf "%.2f" .Price}}</td></tr>
      {{end}}
      </tbody>
    </table>
  </div>

  <div class="section">
    <h2>Lab Test Report</h2>
    <table>
      <thead><tr><th>Test</th><th>Price</th></tr></thead>
      <tbody>
      {{range .LabTests}}
        <tr><td>{{.TestName}}</td><td>{{printf "%.2f" .Price}}</td></tr>
      {{end}}
      </tbody>
    </table>
  </div>

  <div class="section">
    <h2>Billing Report</h2>
    <p><strong>Consultation Fee:</strong> {{printf "%.2f" .Billing.ConsultationFee}}</p>
    <p><strong>Drugs Total:</strong> {{printf "%.2f" .Billing.DrugsTotal}}</p>
    <p><strong>Tests Total:</strong> {{printf "%.2f" .Billing.TestsTotal}}</p>
    <p><strong>Discount:</strong> {{printf "%.2f" .Billing.DiscountAmount}} ({{.Billing.DiscountType}})</p>
    <p><strong>Total:</strong> {{printf "%.2f" .Billing.TotalAmount}}</p>
  </div>
</body>
</html>`

	t, err := template.New("report").Parse(tpl)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return "", err
	}

	return out.String(), nil
}
