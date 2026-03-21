import { ParsedResult } from "../types/domain";

type Props = {
  parsed: ParsedResult;
  setParsed: (value: ParsedResult) => void;
  enabled: boolean;
};

export default function ParsedResultsEditor({ parsed, setParsed, enabled }: Props) {
  function updateDrug(index: number, field: "name" | "dosage" | "frequency", value: string) {
    const next = structuredClone(parsed);
    next.drugs[index][field] = value;
    setParsed(next);
  }

  function addDrug() {
    setParsed({ ...parsed, drugs: [...parsed.drugs, { name: "", dosage: "", frequency: "" }] });
  }

  function removeDrug(index: number) {
    const next = parsed.drugs.filter((_, i) => i !== index);
    setParsed({ ...parsed, drugs: next });
  }

  function updateLab(index: number, value: string) {
    const next = [...parsed.lab_tests];
    next[index] = value;
    setParsed({ ...parsed, lab_tests: next });
  }

  function addLab() {
    setParsed({ ...parsed, lab_tests: [...parsed.lab_tests, ""] });
  }

  function removeLab(index: number) {
    setParsed({ ...parsed, lab_tests: parsed.lab_tests.filter((_, i) => i !== index) });
  }

  return (
    <div className="panel">
      <h2>Parsed Results Editor (Mandatory Verification)</h2>
      {!enabled && <p>Run Process with AI first, then verify and edit all extracted values.</p>}

      <h3>Drugs</h3>
      {parsed.drugs.map((d, idx) => (
        <div key={idx} style={{ display: "grid", gridTemplateColumns: "2fr 1fr 1fr auto", gap: 8, marginBottom: 8 }}>
          <input disabled={!enabled} value={d.name} onChange={(e) => updateDrug(idx, "name", e.target.value)} placeholder="Drug name" />
          <input disabled={!enabled} value={d.dosage} onChange={(e) => updateDrug(idx, "dosage", e.target.value)} placeholder="Dosage" />
          <input disabled={!enabled} value={d.frequency} onChange={(e) => updateDrug(idx, "frequency", e.target.value)} placeholder="Frequency" />
          <button type="button" disabled={!enabled} onClick={() => removeDrug(idx)}>Remove</button>
        </div>
      ))}
      <button type="button" disabled={!enabled} onClick={addDrug}>Add Drug</button>

      <h3>Lab Tests</h3>
      {parsed.lab_tests.map((t, idx) => (
        <div key={idx} style={{ display: "grid", gridTemplateColumns: "1fr auto", gap: 8, marginBottom: 8 }}>
          <input disabled={!enabled} value={t} onChange={(e) => updateLab(idx, e.target.value)} placeholder="Test name" />
          <button type="button" disabled={!enabled} onClick={() => removeLab(idx)}>Remove</button>
        </div>
      ))}
      <button type="button" disabled={!enabled} onClick={addLab}>Add Lab Test</button>

      <h3>Observations</h3>
      <textarea
        rows={5}
        disabled={!enabled}
        value={parsed.observations}
        onChange={(e) => setParsed({ ...parsed, observations: e.target.value })}
      />
    </div>
  );
}
