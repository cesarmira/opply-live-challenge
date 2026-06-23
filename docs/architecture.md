# Architecture Map — Ingredient Substitution API

A wide (not deep) map of the system: what exists today vs. what the
[ROADMAP](../ROADMAP.md) targets. This is an interview-stage POC, so several
layers are intentionally empty today — that contrast is itself the story.

## System at a glance

A single stateless Go HTTP service. One binary, one read endpoint + health,
in-memory data, fully synchronous. The one deliberate design seam is the
`suggest.Suggester` interface, which lets the data source be swapped
(static → database → LLM) without touching the transport layer.

```
client ──GET /suggest?ingredient=butter──▶ ServeMux ──▶ suggestHandler
                                                            │
                                              Suggester.Suggest(ingredient)
                                                            │
                                                Static in-memory map lookup
                                                            │
                                   200 {ingredient, alternatives} | 400 | 404 | 500
```

## Layers

| Layer | Today (built) | Target (ROADMAP) |
|---|---|---|
| **App / entrypoint** | `cmd/server/main.go` — wires deps, reads `PORT`, starts `net/http`. Single binary `bin/server`. | Containerize (Dockerfile); deploy as Lambda. |
| **API / transport** | `internal/api/handler.go` — `ServeMux`: `GET /suggest`, `GET /healthz`. Errors 400/404/500. Depends only on the `Suggester` interface. | Request validation + richer errors; auth, rate limiting, versioning. |
| **Domain / services** | `internal/suggest/suggester.go` (interface, `Alternative`, `ErrNotFound`) + `static.go` (hardcoded map: butter, milk, egg, sugar). | DB-backed and LLM-backed `Suggester`; context-aware (cuisine/allergy); ranking by similarity/availability/cost. |
| **Data** | Hardcoded Go map compiled into the binary. No store. | Real database backend; expanded dataset with dietary tags. |
| **Data pipelines** | **None.** Seed data is source code. | Implied by dataset expansion + ranking → an ingest/enrichment path. |
| **Async jobs** | **None.** Fully synchronous request/response. | Candidates: LLM result caching/precompute, dataset refresh. |
| **UX** | **None.** API-only; `curl` + JSON. README is the human surface. | Not specified — API-first product. |

## Cross-cutting

- **Build/CI:** `Makefile` (build/run/test/smoke/fmt/tidy); `scripts/smoke.sh`
  boots the live server and asserts a known request. CI (GitHub Actions running
  test + smoke) planned.
- **Deploy:** planned Lambda + GitHub CI/CD pipeline.
- **Config:** `PORT` env today; flags / log-level planned.
- **Observability:** none today; structured logging + request IDs planned.
- **Domain rule ([AGENTS.md](../AGENTS.md)):** food & beverages only;
  alternatives must be same-category.

## Talking points

1. The `Suggester` seam is the whole architecture bet — transport never knows
   the data source.
2. Today's "data layer" is a compiled map; the interesting jump is
   static → database → LLM.
3. Pipelines / async / UX are intentionally absent at POC; each maps to a
   concrete ROADMAP item.
4. Verification is a smoke test against the live process, not unit tests (POC
   rule in AGENTS.md).
