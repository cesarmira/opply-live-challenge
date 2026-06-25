# LLM Suggester — Design Spec

**Date:** 2026-06-25
**Status:** Approved

## Goal

Replace the current random-list `Suggester` stub with a real LLM-backed implementation
that calls opencode.ai (DeepSeek v4 Flash) to generate context-aware ingredient alternatives.

## Approach

Approach A — minimal LLM Suggester. One new file, no new dependencies, zero changes to
the HTTP layer or the `Suggester` interface.

## Components & data flow

```
main.go
  reads OPENCODE_API_KEY, OPENCODE_BASE_URL (log.Fatalf if either empty)
  constructs suggest.NewLLM(baseURL, apiKey, "deepseek-v4-flash")
  passes it to api.NewMux(...)  ← unchanged

GET /suggest?ingredient=butter
  → suggestHandler (unchanged)
  → LLM.Suggest("butter")
      POST $OPENCODE_BASE_URL/v1/chat/completions
      body: {model, messages: [system + user]}
  → parse choices[0].message.content as []Alternative
  → return to handler → 200 JSON response
```

### New file: `internal/suggest/llm.go`

```go
type LLM struct {
    baseURL string
    apiKey  string
    model   string
    client  *http.Client  // 10-second timeout
}

func NewLLM(baseURL, apiKey, model string) *LLM
func (l *LLM) Suggest(ingredient string) ([]Alternative, error)
```

Dependencies: stdlib only (`net/http`, `encoding/json`, `strings`).

### `main.go` changes

- Read `OPENCODE_API_KEY` and `OPENCODE_BASE_URL`; `log.Fatalf` if either is empty.
- Replace `suggest.NewList()` with `suggest.NewLLM(baseURL, apiKey, "deepseek-v4-flash")`.
- **Note:** `"deepseek-v4-flash"` is the assumed model ID; verify the exact string against opencode.ai's model list before shipping.

## Prompt design

**System:**
```
You are a culinary expert. Given a food or beverage ingredient, return 3–5
alternatives from the same category (food for food, beverage for beverage).
Respond ONLY with a valid JSON array. Each element: {"name": "...", "notes": "..."}.
Never include non-edible items. No markdown, no explanation.
```

**User:** `Suggest alternatives for: {ingredient}`

## Response parsing

Unmarshal `choices[0].message.content` into `[]Alternative`. Strip leading/trailing
whitespace and backtick code fences before unmarshalling as a safety measure.

## Error handling

| Scenario | Behaviour |
|---|---|
| Missing env var at startup | `log.Fatalf` — server does not start |
| Network / non-2xx from API | `Suggest` returns error → handler 500 |
| Unparseable JSON response | `Suggest` returns error → handler 500 |
| LLM call hangs | 10-second `http.Client` timeout → error → 500 |

No silent fallbacks. The caller always sees a real failure.

## Files changed

| File | Change |
|---|---|
| `internal/suggest/llm.go` | **New** — `LLM` struct + `NewLLM` + `Suggest` |
| `cmd/server/main.go` | Read env vars; wire `NewLLM` instead of `NewList` |
| `internal/suggest/list.go` | Kept (not deleted — still referenced nowhere after swap, can be cleaned later) |
