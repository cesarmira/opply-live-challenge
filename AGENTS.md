# AGENTS.md

Working agreement for any AI agent (and humans) contributing to this project.

## What this project is

An HTTP API that suggests **alternative ingredients** for a given ingredient.
Input is JSON, output is JSON.

```
POST /suggest   {"ingredient": "butter"}  ->  {"ingredient": "butter", "alternatives": [...]}
```

## Rules

1. **Language: Go.** All application code is written in Go using the standard
   library by default. Add third-party dependencies only when they earn their place.

2. **Keep units small and isolated.** Each package and file has one clear
   purpose and communicates through well-defined interfaces. The HTTP layer
   (`internal/api`) depends only on the `suggest.Suggester` interface — never on
   a concrete implementation. If a file grows large or hard to follow, that is a
   signal to split it.

3. **No unit tests for now — this is a POC.** Do not add or maintain unit tests
   while we are proving the concept. The smoke test is the only verification gate
   (see rule 5). Revisit this once the project moves past POC.

4. **Commit and push to `main` after each successful change.** A change is
   "successful" once it builds (`make build`) and the smoke test passes. Then
   commit with a clear message and push to `main`. Small, frequent commits over
   large ones.

5. **Run a smoke test after each push.** After pushing, run `make smoke` to
   verify the running API still answers a known request correctly. If the smoke
   test fails, fix it before starting new work.

## Standard loop for every change

```
make fmt        # format
make build      # must compile
git commit ...  # clear message
git push        # to main
make smoke      # verify the live API
```

## Layout

```
cmd/server/        entrypoint: wires dependencies, starts the HTTP server
internal/suggest/  domain: Suggester interface + implementations
internal/api/      transport: JSON in/out over HTTP, no domain logic
scripts/           helper scripts (smoke test)
```
