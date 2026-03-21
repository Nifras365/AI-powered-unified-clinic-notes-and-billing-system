package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"clinicnotes/backend/internal/dto"
)

const basePrompt = `You are a medical information extraction system.

Extract structured medical data from the clinical text below.

Return ONLY valid JSON in this exact format:
{
  "drugs": [
    { "name": "", "dosage": "", "frequency": "" }
  ],
  "lab_tests": [],
  "observations": ""
}

STRICT RULES:
- Do NOT include explanations
- Do NOT include markdown
- Only output JSON
- If no items exist, return empty arrays
- Ensure valid JSON syntax
- Normalize names (drug names, test names)

Clinical Text:
%s`

const strictRetryPrompt = `Output invalid previously. Respond with ONLY strict JSON. No markdown, no prose.
Use exact schema and valid syntax.
Clinical Text:
%s`

type AIService interface {
	ParseMedicalText(ctx context.Context, input string) (dto.ParsedResult, error)
}

type AIOrchestrator struct {
	primary  AIService
	fallback AIService
}

func NewAIOrchestrator(primary AIService, fallback AIService) *AIOrchestrator {
	return &AIOrchestrator{primary: primary, fallback: fallback}
}

func (o *AIOrchestrator) ParseMedicalText(ctx context.Context, input string) (dto.ParsedResult, error) {
	result, err := o.primary.ParseMedicalText(ctx, input)
	if err == nil {
		return result, nil
	}
	if o.fallback == nil {
		return dto.ParsedResult{}, err
	}
	fallbackResult, fallbackErr := o.fallback.ParseMedicalText(ctx, input)
	if fallbackErr == nil {
		return fallbackResult, nil
	}

	return dto.ParsedResult{}, fmt.Errorf("primary provider failed: %v; fallback provider failed: %v", err, fallbackErr)
}

func parseAndValidateJSON(raw string) (dto.ParsedResult, error) {
	cleaned := cleanAIOutput(raw)
	jsonPayload, err := extractJSONObject(cleaned)
	if err != nil {
		return dto.ParsedResult{}, err
	}

	var parsed dto.ParsedResult
	if err := json.Unmarshal([]byte(jsonPayload), &parsed); err != nil {
		return dto.ParsedResult{}, fmt.Errorf("invalid json: %w", err)
	}

	if err := validateParsedResult(parsed); err != nil {
		return dto.ParsedResult{}, err
	}

	return parsed, nil
}

func cleanAIOutput(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	return strings.TrimSpace(raw)
}

func extractJSONObject(input string) (string, error) {
	re := regexp.MustCompile(`(?s)\{.*\}`)
	match := re.FindString(input)
	if match == "" {
		return "", errors.New("no json object found")
	}
	return match, nil
}

func validateParsedResult(parsed dto.ParsedResult) error {
	if parsed.Drugs == nil {
		parsed.Drugs = []dto.Drug{}
	}
	if parsed.LabTests == nil {
		parsed.LabTests = []string{}
	}
	for _, d := range parsed.Drugs {
		if strings.TrimSpace(d.Name) == "" {
			return errors.New("drug name cannot be empty")
		}
	}
	return nil
}
