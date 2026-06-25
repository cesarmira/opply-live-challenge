package suggest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// LLM is a Suggester backed by an OpenAI-compatible chat completion API.
type LLM struct {
	baseURL string
	apiKey  string
	model   string
	client  *http.Client
}

// NewLLM returns an LLM wired to the given base URL, API key, and model.
func NewLLM(baseURL, apiKey, model string) *LLM {
	return &LLM{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

const systemPrompt = `You are a culinary expert. Given a food or beverage ingredient, return 3–5 alternatives from the same category (food for food, beverage for beverage). Respond ONLY with a valid JSON array. Each element: {"name": "...", "notes": "..."}. Never include non-edible items. No markdown, no explanation.`

// Suggest calls the LLM and returns ingredient alternatives.
func (l *LLM) Suggest(ingredient string) ([]Alternative, error) {
	payload := chatRequest{
		Model: l.model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: "Suggest alternatives for: " + ingredient},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, l.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.apiKey)

	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call LLM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("LLM returned status %d", resp.StatusCode)
	}

	var chat chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	if len(chat.Choices) == 0 {
		return nil, ErrNotFound
	}

	content := stripFences(chat.Choices[0].Message.Content)

	var alts []Alternative
	if err := json.Unmarshal([]byte(content), &alts); err != nil {
		return nil, fmt.Errorf("parse alternatives: %w", err)
	}
	if len(alts) == 0 {
		return nil, ErrNotFound
	}
	return alts, nil
}

// stripFences removes whitespace and markdown code fences LLMs occasionally add.
func stripFences(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "```") {
		if i := strings.Index(s, "\n"); i != -1 {
			s = s[i+1:]
		}
		if i := strings.LastIndex(s, "```"); i != -1 {
			s = s[:i]
		}
		s = strings.TrimSpace(s)
	}
	return s
}
