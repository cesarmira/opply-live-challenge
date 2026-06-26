# ROADMAP

## Today pre-interview (boilerplate)

- [x] Project layout: `cmd/server`, `internal/api`, `internal/suggest`
- [x] `Suggester` interface + static in-memory stub implementation
- [x] `POST /suggest` endpoint: JSON in → JSON out
- [x] `GET /healthz` liveness endpoint
- [x] Unit tests for the domain and the HTTP handler
- [x] Makefile: `build`, `run`, `test`, `smoke`, `fmt`, `tidy`, `clean`
- [x] Smoke test (`make smoke`) that exercises the live API
- [x] Working agreement (`AGENTS.md`) and `CLAUDE.md` pointer

## Today live challenge

- [ ] Request a new ingredient to an LLM (https://opencode.ai/)
- [ ] Deploy application as _probably_ lambda locally
- [ ] Deploy the application in a pipeline (GitHub CI/CD)
- [ ] Talk with interviewers about new items in the ROADMAP
- [ ] Receive/fix feedback
- [ ] Make a good impression

## Next

- [ ] Expand the substitution dataset (more ingredients, dietary tags)
- [ ] Support quantity/ratio hints in responses (e.g. "use 3/4 the amount")
- [ ] Request validation and richer error responses
- [ ] Structured logging and request IDs
- [ ] Configurable via environment / flags (port, log level)
- [ ] Containerize (Dockerfile) and add CI running `make test` + `make smoke`
    
## Future

- [ ] Swap the static stub for a real data backend (database)
- [ ] LLM-backed `Suggester` for open-ended / unknown ingredients
- [ ] Context-aware suggestions (recipe type, allergies, cuisine)
- [ ] Rank alternatives by similarity / availability / cost
- [ ] Public API: auth, rate limiting, versioning
