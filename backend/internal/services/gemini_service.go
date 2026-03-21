package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"clinicnotes/backend/internal/dto"
)

type GeminiService struct {
	apiKey string
	model  string
	client *http.Client
}

type geminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func NewGeminiService(apiKey, model string) *GeminiService {
	return &GeminiService{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 25 * time.Second},
	}
}

func (s *GeminiService) ParseMedicalText(ctx context.Context, input string) (dto.ParsedResult, error) {
	if s.apiKey == "" {
		return dto.ParsedResult{}, fmt.Errorf("gemini api key is empty")
	}

	prompts := []string{
		fmt.Sprintf(basePrompt, input),
		fmt.Sprintf(strictRetryPrompt, input),
		fmt.Sprintf(strictRetryPrompt, input),
	}

	var lastErr error
	for _, prompt := range prompts {
		raw, err := s.callGemini(ctx, prompt)
		if err != nil {
			lastErr = err
			continue
		}

		parsed, err := parseAndValidateJSON(raw)
		if err != nil {
			lastErr = err
			continue
		}

		return parsed, nil
	}

	return dto.ParsedResult{}, fmt.Errorf("gemini parse failed after retries: %w", lastErr)
}

func (s *GeminiService) callGemini(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", s.model, s.apiKey)

	payload := geminiRequest{}
	content := struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{}
	part := struct {
		Text string `json:"text"`
	}{Text: prompt}
	content.Parts = []struct {
		Text string `json:"text"`
	}{part}
	payload.Contents = []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}{content}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("gemini api error (%d): %s", resp.StatusCode, string(respBody))
	}

	var parsedResp geminiResponse
	if err := json.Unmarshal(respBody, &parsedResp); err != nil {
		return "", err
	}

	if len(parsedResp.Candidates) == 0 || len(parsedResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned no content")
	}

	return parsedResp.Candidates[0].Content.Parts[0].Text, nil
}
