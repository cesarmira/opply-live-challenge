# LLM Suggester Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the random-list `Suggester` stub with an LLM-backed implementation that calls opencode.ai (DeepSeek v4 Flash) to return real, context-aware ingredient alternatives.

**Architecture:** Add one new file `internal/suggest/llm.go` that implements the existing `Suggester` interface by calling an OpenAI-compatible chat completion endpoint. Update `main.go` to read two env vars and wire the new implementation. The HTTP layer (`internal/api`) is untouched.

**Tech Stack:** Go 1.22 stdlib only — `net/http`, `encoding/json`, `bytes`, `strings`, `fmt`, `time`.

## Global Constraints

- Go only; no third-party dependencies (stdlib exclusively)
- No unit tests — POC; smoke test is the only verification gate
- Every file change must pass `make fmt && make build` before committing
- After the final push, run `make smoke` to verify the live API
- Commit message style: imperative, lower-case, one sentence
- `OPENCODE_API_KEY` and `OPENCODE_BASE_URL` must be set in the shell before running the smoke test
- Model ID to use: `"deepseek-v4-flash"` — verify this matches the exact string opencode.ai expects

---

### Task 1: Create the LLM Suggester

**Files:**
- Create: `internal/suggest/llm.go`

**Interfaces:**
- Produces: `suggest.NewLLM(baseURL, apiKey, model string) *LLM` and `(*LLM).Suggest(ingredient string) ([]Alternative, error)` — satisfies `suggest.Suggester`

- [ ] **Step 1: Create `internal/suggest/llm.go`**

Write the file with the exact content below:

```go
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
		baseURL: baseURL,
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
```

- [ ] **Step 2: Format and build**

```bash
make fmt && make build
```

Expected: no output errors, `bin/server` is produced.

- [ ] **Step 3: Commit**

```bash
git add internal/suggest/llm.go
git commit -m "add LLM-backed Suggester using opencode.ai / DeepSeek v4 Flash"
```

---

### Task 2: Wire LLM Suggester in main.go and verify end-to-end

**Files:**
- Modify: `cmd/server/main.go`

**Interfaces:**
- Consumes: `suggest.NewLLM(baseURL, apiKey, model string) *LLM` from Task 1

- [ ] **Step 1: Update `cmd/server/main.go`**

Replace the entire file with:

```go
// Command server starts the ingredient-substitution HTTP API.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/cesarmirasanchez/opply-live-challenge/internal/api"
	"github.com/cesarmirasanchez/opply-live-challenge/internal/suggest"
)

func main() {
	apiKey := os.Getenv("OPENCODE_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENCODE_API_KEY is required")
	}
	baseURL := os.Getenv("OPENCODE_BASE_URL")
	if baseURL == "" {
		log.Fatal("OPENCODE_BASE_URL is required")
	}

	addr := ":" + port()
	mux := api.NewMux(suggest.NewLLM(baseURL, apiKey, "deepseek-v4-flash"))

	log.Printf("ingredient-substitution API listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

// port returns the listen port from PORT, defaulting to 8080.
func port() string {
	if p := os.Getenv("PORT"); p != "" {
		return p
	}
	return "8080"
}
```

- [ ] **Step 2: Format and build**

```bash
make fmt && make build
```

Expected: compiles cleanly, `bin/server` updated.

- [ ] **Step 3: Export credentials and run smoke test**

The smoke script starts `bin/server` as a subprocess; it inherits env vars from the shell, so export them before running:

```bash
export OPENCODE_API_KEY=<your-key>
export OPENCODE_BASE_URL=<https://your-opencode-base-url>
make smoke
```

Expected final line: `SMOKE OK`

If the LLM returns alternatives, the response will look like:
```json
{"ingredient":"butter","alternatives":[{"name":"ghee","notes":"..."},{"name":"coconut oil","notes":"..."}]}
```

If the smoke test fails with a 500, check that `OPENCODE_BASE_URL` does not have a trailing slash and that the model ID `"deepseek-v4-flash"` matches the exact string the opencode.ai API expects.

- [ ] **Step 4: Commit and push**

```bash
git add cmd/server/main.go
git commit -m "wire LLM Suggester in main; read OPENCODE_API_KEY and OPENCODE_BASE_URL"
git push
```

- [ ] **Step 5: Re-run smoke after push to confirm live API**

```bash
make smoke
```

Expected: `SMOKE OK`
