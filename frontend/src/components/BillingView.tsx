import { Billing } from "../types/domain";

type Props = {
  billing: Billing | null;
};

function money(value: number | undefined | null) {
  const num = typeof value === "number" ? value : Number(value ?? 0);
  return Number.isFinite(num) ? num.toFixed(2) : "0.00";
}

export default function BillingView({ billing }: Props) {
  return (
    <div className="panel">
      <h2>Billing Summary</h2>
      {!billing && <p>Save a verified consultation to generate billing.</p>}
      {billing && (
        <table>
          <tbody>
            <tr><td>Consultation Fee</td><td>{money(billing.consultation_fee)}</td></tr>
            <tr><td>Drugs Total</td><td>{money(billing.drugs_total)}</td></tr>
            <tr><td>Lab Tests Total</td><td>{money(billing.tests_total)}</td></tr>
            <tr><td>Subtotal</td><td>{money(billing.subtotal)}</td></tr>
            <tr><td>Discount ({billing.discount_type})</td><td>-{money(billing.discount_amount)}</td></tr>
            <tr><td><strong>Total</strong></td><td><strong>{money(billing.total_amount)}</strong></td></tr>
          </tbody>
        </table>
      )}
    </div>
  );
}
