import { ParsedResult } from "../types/domain";

type Props = {
  text: string;
  parsed: ParsedResult;
};

function escapeRegExp(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

export default function HighlightedTextViewer({ text, parsed }: Props) {
  const drugNames = parsed.drugs.map((d) => d.name).filter(Boolean);
  const testNames = parsed.lab_tests.filter(Boolean);

  if (!text.trim()) {
    return (
      <div className="panel">
        <h2>Entity Highlight Viewer</h2>
        <p>Enter clinical notes and parse to see highlighted drugs and lab tests.</p>
      </div>
    );
  }

  let html = text;

  for (const drug of drugNames) {
    const reg = new RegExp(`(${escapeRegExp(drug)})`, "gi");
    html = html.replace(reg, '<mark class="drug">$1</mark>');
  }

  for (const test of testNames) {
    const reg = new RegExp(`(${escapeRegExp(test)})`, "gi");
    html = html.replace(reg, '<mark class="test">$1</mark>');
  }

  return (
    <div className="panel">
      <h2>AI Verification Highlight</h2>
      <p>
        <span className="badge drug">Drug</span>
        <span className="badge test">Lab Test</span>
      </p>
      <div className="highlight-viewer" dangerouslySetInnerHTML={{ __html: html }} />
    </div>
  );
}
