import { useMemo, useState } from "react";
import { getBilling, getReportHTML, parseClinicalText, saveConsultation } from "./api";
import BillingView from "./components/BillingView";
import HighlightedTextViewer from "./components/HighlightedTextViewer";
import InputPanel from "./components/InputPanel";
import ParsedResultsEditor from "./components/ParsedResultsEditor";
import ReportView from "./components/ReportView";
import { Billing, ParsedResult } from "./types/domain";

const emptyParsed: ParsedResult = {
  drugs: [],
  lab_tests: [],
  observations: "",
};

export default function App() {
  const [rawInput, setRawInput] = useState("");
  const [patientName, setPatientName] = useState("");
  const [patientAge, setPatientAge] = useState<number>(35);
  const [parsed, setParsed] = useState<ParsedResult>(emptyParsed);
  const [hasParsed, setHasParsed] = useState(false);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [billing, setBilling] = useState<Billing | null>(null);
  const [reportHTML, setReportHTML] = useState<string>("");
  const [consultationFee, setConsultationFee] = useState(40);
  const [discountType, setDiscountType] = useState<"fixed" | "percentage">("fixed");
  const [discountValue, setDiscountValue] = useState(0);

  const drugPrices = useMemo(
    () => Object.fromEntries(parsed.drugs.map((d) => [d.name, 0])),
    [parsed.drugs]
  );
  const testPrices = useMemo(
    () => Object.fromEntries(parsed.lab_tests.map((t) => [t, 0])),
    [parsed.lab_tests]
  );

  async function handleProcess() {
    setError(null);
    setBilling(null);
    setReportHTML("");

    if (!rawInput.trim()) {
      setError("Please enter clinical notes before processing.");
      return;
    }

    try {
      setLoading(true);
      const result = await parseClinicalText(rawInput);
      setParsed(result);
      setHasParsed(true);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setLoading(false);
    }
  }

  async function handleSave() {
    setError(null);

    if (!hasParsed) {
      setError("Run AI parse first, then verify/edit before saving.");
      return;
    }
    if (!patientName.trim() || patientAge <= 0) {
      setError("Patient name and age are required.");
      return;
    }

    try {
      setSaving(true);
      const saveRes = await saveConsultation({
        patient: { name: patientName.trim(), age: patientAge },
        raw_input: rawInput,
        parsed,
        pricing: {
          consultation_fee: consultationFee,
          discount_type: discountType,
          discount_value: discountValue,
          drug_prices: drugPrices,
          test_prices: testPrices,
        },
      });

      const [billingRes, reportRes] = await Promise.all([
        getBilling(saveRes.id),
        getReportHTML(saveRes.id),
      ]);

      setBilling(billingRes);
      setReportHTML(reportRes);
    } catch (e) {
      setError((e as Error).message);
    } finally {
      setSaving(false);
    }
  }

  return (
    <div className="app-shell">
      <header className="hero">
        <p className="eyebrow">AI-Powered Unified Clinic Notes & Billing</p>
        <h1>Single workspace for clinical dictation, structured extraction, and billing.</h1>
      </header>

      <section className="panel-grid">
        <InputPanel
          rawInput={rawInput}
          setRawInput={setRawInput}
          patientName={patientName}
          setPatientName={setPatientName}
          patientAge={patientAge}
          setPatientAge={setPatientAge}
          onProcess={handleProcess}
          processing={loading}
        />

        <HighlightedTextViewer text={rawInput} parsed={parsed} />
      </section>

      <section className="panel-grid single">
        <ParsedResultsEditor parsed={parsed} setParsed={setParsed} enabled={hasParsed} />
      </section>

      <section className="billing-controls">
        <h3>Billing Controls</h3>
        <div className="billing-fields">
          <label>
            Consultation Fee
            <input type="number" value={consultationFee} onChange={(e) => setConsultationFee(Number(e.target.value))} />
          </label>
          <label>
            Discount Type
            <select value={discountType} onChange={(e) => setDiscountType(e.target.value as "fixed" | "percentage")}>
              <option value="fixed">Fixed</option>
              <option value="percentage">Percentage</option>
            </select>
          </label>
          <label>
            Discount Value
            <input type="number" value={discountValue} onChange={(e) => setDiscountValue(Number(e.target.value))} />
          </label>
        </div>
        <button className="save-btn" disabled={saving || !hasParsed} onClick={handleSave}>
          {saving ? "Saving..." : "Confirm & Save Consultation"}
        </button>
      </section>

      {error && <p className="error-banner">{error}</p>}

      <section className="panel-grid">
        <BillingView billing={billing} />
        <ReportView reportHTML={reportHTML} />
      </section>
    </div>
  );
}
