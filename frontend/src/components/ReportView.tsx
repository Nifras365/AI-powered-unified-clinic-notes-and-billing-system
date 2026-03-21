import html2pdf from "html2pdf.js";

type Props = {
  reportHTML: string;
};

export default function ReportView({ reportHTML }: Props) {
  function handlePrint() {
    if (!reportHTML) return;

    const printWindow = window.open("", "_blank", "width=900,height=700");
    if (!printWindow) {
      window.alert("Pop-up blocked. Please allow pop-ups to print the report.");
      return;
    }

    printWindow.document.open();
    printWindow.document.write(reportHTML);
    printWindow.document.close();
    printWindow.focus();
    printWindow.print();
  }

  async function handleDownloadPdf() {
    if (!reportHTML) return;

    const container = document.createElement("div");
    container.style.padding = "16px";
    container.style.background = "#ffffff";
    container.innerHTML = reportHTML;

    const options = {
      margin: 10,
      filename: `clinic-report-${new Date().toISOString().slice(0, 10)}.pdf`,
      image: { type: "jpeg", quality: 0.98 },
      html2canvas: { scale: 2, useCORS: true },
      jsPDF: { unit: "mm", format: "a4", orientation: "portrait" },
    };

    await html2pdf().set(options).from(container).save();
  }

  return (
    <div className="panel">
      <h2>Report Preview (Printable)</h2>
      <div className="report-actions">
        <button type="button" onClick={handlePrint} disabled={!reportHTML}>
          Print
        </button>
        <button type="button" className="save-btn" onClick={handleDownloadPdf} disabled={!reportHTML}>
          Download PDF
        </button>
      </div>
      {!reportHTML && <p>Report will appear after consultation save.</p>}
      {reportHTML && <iframe title="report" className="report-frame" srcDoc={reportHTML} />}
    </div>
  );
}
