export type Drug = {
  name: string;
  dosage: string;
  frequency: string;
};

export type ParsedResult = {
  drugs: Drug[];
  lab_tests: string[];
  observations: string;
};

export type Billing = {
  id: string;
  consultation_id: string;
  consultation_fee: number;
  drugs_total: number;
  tests_total: number;
  subtotal: number;
  discount_type: string;
  discount_value: number;
  discount_amount: number;
  total_amount: number;
};

export type ConsultationCreatePayload = {
  patient: { name: string; age: number };
  raw_input: string;
  parsed: ParsedResult;
  pricing: {
    consultation_fee: number;
    discount_type: "fixed" | "percentage";
    discount_value: number;
    drug_prices: Record<string, number>;
    test_prices: Record<string, number>;
  };
};
