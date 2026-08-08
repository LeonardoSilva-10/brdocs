# Changelog

Todas as mudanças notáveis deste projeto serão documentadas aqui.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [Unreleased]

### Added
- Geração aleatória de CPF e CNPJ válidos (`cpf.Generate`, `cnpj.Generate`, `brdocs.GenerateCPF`/`GenerateCNPJ`)
- Cache de consultas nas APIs via `api.WithCacheTTL`
- Resultado tipado para a API cpfcnpj.com.br (`api.CPFCNPJResult`) e método `Client.ConsultaCPFCNPJ`

### Fixed
- Exemplo `ExampleState` usava CPF inválido e derrubava `go test ./...`
- `CNPJResult.IsActive` passou a considerar apenas situação ATIVA (código 2); códigos 1 (NULA), 3 (SUSPENSA), 4 (INAPTA) e 8 (BAIXADA) agora retornam `false`
- Token OAuth do Serpro agora é reutilizado entre chamadas do mesmo `Client`
- Removidos os sentinelas mortos `ErrInvalidDigits` dos pacotes `cpf` e `cnpj`

### Removed
- Pacote interno `internal/httpclient` (código morto, não utilizado pelo `api`)

## [1.0.0] - 2026-08-07

### Added
- Lançamento inicial do pacote brdocs
