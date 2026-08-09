# agent.md — Guide for agents (and contributors) in brdocs

This document guides anyone — human or AI agent — who will touch this
repository. Read it before editing code, tests, or documentation.

## 1. What this project is

`brdocs` is a **Go library** for validation, generation, and querying of
Brazilian CPF and CNPJ numbers. Nothing here is an executable application: it
is a publishable Go module (`github.com/LeonardoSilva-10/brdocs`), so the
**public API matters** to whoever consumes it.

- Local validation: no external dependencies (stdlib only) in the `cpf` and `cnpj` packages.
- Queries: `api` talks to Serpro/ConectaGov, Receita Federal (scraping), and
  third-party services (cpfcnpj.com.br, cpf-brasil.org).
- The only `go.mod` dependency is `testify` (used **only in tests**).

## 2. Actual structure (do not invent folders that do not exist)

```
brdocs/
├── cpf/cpf.go              # CPF validation, formatting, masking, state and generation
├── cnpj/cnpj.go            # CNPJ validation, formatting, masking, branch and generation
├── api/
│   ├── client.go           # Client, options, cache and sanitizeDocument
│   ├── serpro.go           # Serpro/ConectaGov (CNPJ) — OAuth2
│   ├── receitafederal.go   # Receita Federal public query (CPF) — scraping
│   ├── thirdparty.go       # cpfcnpj.com.br and cpf-brasil.org
│   └── types.go            # Result structs and helpers (IsActive, IsRegular)
├── brdocs.go               # Convenience re-exports (cpf/cnpj)
├── e2e/e2e_test.go         # End-to-end tests (some require credentials/env)
└── .github/workflows/ci.yml
```

> ⚠️ **Do not** add `internal/`, `mocks/`, etc. folders without a real need.
> The old `internal/httpclient` was removed because it was dead code. Each
> package inside `api/` is a `.go` file in the same `api` package, not a subfolder.

## 3. Code conventions

- **Idiomatic Go**: follow `Effective Go` and `gofmt`. Run `gofmt -w` before committing.
- **Every exported function/type/var has a docstring** starting with the symbol name.
- **Errors**: use `errors.New` for sentinels (`ErrInvalidLength`, etc.) and
  `fmt.Errorf("pkg: ...", %w)` to add context. Prefix the message with the
  package (`api:`, `serpro:`, `cpfcnpj:`).
- **Context**: every API call receives `context.Context` as its first argument.
- **Tests**: use `testify` (`assert`/`require`). Prefer table-driven tests.
- **Examples**: `ExampleXxx` runs **as a test** — the `// Output:` block is
  verified. Never use an invalid CPF/CNPJ in an example without checking the output.
- **Dependencies**: do not add external dependencies to the validation packages
  (`cpf`/`cnpj`). Tests/API may use what already exists (`testify`, stdlib).

## 4. Domain knowledge (do not get these wrong)

- **CPF**: 11 digits, 2 check digits mod-11 (weights 10..2 and 11..2). All-equal
  digits (e.g. `111.111.111-11`) are invalid.
- **CNPJ**: 14 digits, 2 check digits mod-11 (weights 5..2 and 6..2). All-equal
  digits make it invalid. Branch/headquarters: positions 9–12 (`0001` = headquarters).
- **CPF state (UF)**: the **9th digit** indicates the issuing state. A digit may
  map to several states (returns `[]string`). E.g. digit 2 → AC/AM/AP/PA/RO/RR.
- **CNPJ registration status (Receita Federal)** — official table:
  `1`=NULA, `2`=ATIVA, `3`=SUSPENSA, `4`=INAPTA, `8`=BAIXADA.
  So `CNPJResult.IsActive()` ONLY considers code `2` (or a description containing
  "ATIVA"). **Never** treat `1` as active.
- **CPF status**: textual values (REGULAR / IRREGULAR / CANCELADA / NULA…).

## 5. Query APIs — rules when touching

- **Cache**: implemented in `api` via `WithCacheTTL` (disabled by default).
  Keys are `provider:document` (Receita Federal also includes the birth date).
  Only **successful responses** are cached. When adding a provider, use
  `cacheGet`/`cacheSet` with your own prefix.
- **Serpro OAuth token**: lives on the `serproClient` **persisted on the `Client`**
  (`c.serpro`) so it is reused across calls. Do not recreate the client per call.
- **Sanitization**: use `sanitizeDocument` (removes non-digits) instead of
  repeating `strings.ReplaceAll` in each file.
- **Third-party response (cpfcnpj)**: typed (`CPFCNPJResult`), not `interface{}`.
- **doRequest**: only decodes 200 OK; status errors return a generic error. Do
  not cache errors.

## 6. Fragile / known areas

- **Receita Federal (CPF)**: it is **web scraping** of a public page. The site
  has CAPTCHA and changes its HTML structure; the parsing regexes may break at
  any time. That is why the E2E test is behind `RUN_RF_INTEGRATION=1`.
- **CPF validation via RF** depends on the birth date and an internal check
  digit. Do not assume that just hitting the URL works — it is fragile.
- **`go test ./...` must stay green ALWAYS**, including the examples
  (`ExampleState` already broke the pipeline by using an invalid CPF).

## 7. Commands

```bash
go test ./...                 # all tests (must pass)
go test -race ./...           # with race detector
go test -cover ./...          # per-package coverage
go vet ./...                  # static analysis
gofmt -w $(find . -name '*.go')   # formatting
go test -bench=. -benchmem ./...  # benchmarks

# E2E (require credentials/environment; CI skips without them)
SERPRO_CONSUMER_KEY=x SERPRO_CONSUMER_SECRET=x go test ./e2e/... -run TestSerproIntegration
CPFCNPJ_API_KEY=x go test ./e2e/... -run TestCPFCNPJAPIIntegration
RUN_RF_INTEGRATION=1 go test ./e2e/... -run TestReceitaFederalIntegration
```

Alternatively there is a `Makefile` with `make test`, `make coverage`, `make lint`,
`make bench`, `make e2e`.

## 8. Mandatory checklist for any change

1. Understand the impact on the **public API** (exported functions/structs).
2. Run `gofmt` and `go vet ./...`.
3. Run `go test ./...` — **including the `ExampleXxx` functions** — and see it all green.
4. If you change validator behavior, run `go test -cover ./...` and avoid
   meaningful coverage regressions.
5. Update `README.md` (features, examples, API table, structure) and
   `CHANGELOG.md` when observable behavior changes.
6. Keep `agent.md` and `README` **in sync with the real code** (folders,
   function names, naming conventions). E.g. on the root it is
   `brdocs.IsValidCPF`; in the subpackages it is `cpf.Valid`/`cnpj.Valid`; the
   `api` package uses `QueryCNPJ`/`QueryCPF`/`QueryCPFCNPJ` — do not invent
   `*.IsValid` or `*.ConsultaCNPJ`.
7. Never embed credentials/variables in code — read them from the environment
   in E2E tests and skip (`t.Skip`) when absent.

## 9. Privacy and legal notice

CPF data is protected by the **LGPD** (Brazilian General Data Protection Law).
Use valid test CPFs/CNPJs and never collect or persist personal data. Official
APIs have their own terms of service (Serpro/Receita Federal). Document any
CAPTCHA/credential limits or requirements.