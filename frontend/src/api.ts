import { Billing, ConsultationCreatePayload, ParsedResult } from "./types/domain";

// in api backend host 8080
const API_BASE = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE}${path}`, {
    headers: { "Content-Type": "application/json", ...(init?.headers || {}) },
    ...init,
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: "Request failed" }));
    throw new Error(body.error || "Request failed");
  }

  const contentType = response.headers.get("content-type") || "";
  if (contentType.includes("application/json")) {
    return (await response.json()) as T;
  }

  return (await response.text()) as T;
}

export async function parseClinicalText(rawInput: string): Promise<ParsedResult> {
  const res = await request<{ parsed: ParsedResult }>("/api/parse", {
    method: "POST",
    body: JSON.stringify({ raw_input: rawInput }),
  });
  return res.parsed;
}

export async function saveConsultation(payload: ConsultationCreatePayload): Promise<{ id: string }> {
  return request<{ id: string }>("/api/consultations", {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function getBilling(consultationId: string): Promise<Billing> {
  return request<Billing>(`/api/billing/${consultationId}`);
}

export async function getReportHTML(consultationId: string): Promise<string> {
  return request<string>(`/api/reports/${consultationId}`);
}
