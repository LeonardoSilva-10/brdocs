# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added
- Random generation of valid CPF and CNPJ numbers (`cpf.Generate`, `cnpj.Generate`, `brdocs.GenerateCPF`/`GenerateCNPJ`)
- Query result caching via `api.WithCacheTTL`
- Typed result for the cpfcnpj.com.br API (`api.CPFCNPJResult`) and the `Client.QueryCPFCNPJ` method

### Changed
- **Breaking:** the `api` package public identifiers were translated to English (`Client.ConsultaCNPJ` → `QueryCNPJ`, `SituacaoCadastral` → `RegistrationStatus`, `Codigo` → `Code`, etc.). JSON tags were kept unchanged to preserve compatibility with the upstream providers.

### Fixed
- `ExampleState` used an invalid CPF and broke `go test ./...`
- `CNPJResult.IsActive` now only considers ATIVA (code 2); codes 1 (NULA), 3 (SUSPENSA), 4 (INAPTA), and 8 (BAIXADA) now return `false`
- The Serpro OAuth token is now reused across calls of the same `Client`
- Removed the dead `ErrInvalidDigits` sentinels from the `cpf` and `cnpj` packages

### Removed
- The internal `internal/httpclient` package (dead code, not used by `api`)

## [1.0.0] - 2026-08-07

### Added
- Initial release of the brdocs package