# brdocs

[![Go Reference](https://pkg.go.dev/badge/github.com/LeonardoSilva-10/brdocs.svg)](https://pkg.go.dev/github.com/LeonardoSilva-10/brdocs)
[![CI](https://github.com/LeonardoSilva-10/brdocs/actions/workflows/ci.yml/badge.svg)](https://github.com/LeonardoSilva-10/brdocs/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/LeonardoSilva-10/brdocs)](https://goreportcard.com/report/github.com/LeonardoSilva-10/brdocs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Go package to validate and query Brazilian CPF and CNPJ numbers.

## ✨ Features

- ✅ **Mathematical validation** of CPF and CNPJ (official mod-11 algorithm)
- 🔍 **Queries through APIs** from the Brazilian government and third-party services
- 🛡️ **Automatic sanitization** (removes dots, dashes, slashes)
- 🎲 **Generation** of valid CPF and CNPJ (`crypto/rand`)
- 📦 **No external dependencies** for local validation
- 🧪 **Unit and E2E tests** with high validator coverage (~90%)
- ⚡ **Built-in cache** for repeated queries (`api.WithCacheTTL`)
- 🌐 **Native context support** (`context.Context`) in all API calls

## 📦 Installation

```bash
go get github.com/LeonardoSilva-10/brdocs
```

## 🚀 Quick Start

### Local validation (no internet)

```go
package main

import (
    "fmt"

    "github.com/LeonardoSilva-10/brdocs/cpf"
    "github.com/LeonardoSilva-10/brdocs/cnpj"
)

func main() {
    // CPF
    fmt.Println(cpf.Valid("529.982.247-25")) // true
    fmt.Println(cpf.Valid("111.111.111-11")) // false (repeated digits)
    fmt.Println(cpf.Valid("123.456.789-09")) // false (invalid check digits)

    // CNPJ
    fmt.Println(cnpj.Valid("11.222.333/0001-81")) // true
    fmt.Println(cnpj.Valid("00.000.000/0000-00")) // false
}
```

### API query

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/LeonardoSilva-10/brdocs/api"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client := api.NewClient(
        api.WithSerproCredentials("consumerKey", "consumerSecret"),
        api.WithCacheTTL(5*time.Minute), // cache query results for 5 minutes
    )

    // Query a CNPJ through the Receita Federal (Serpro)
    result, err := client.QueryCNPJ(ctx, "34238864000168")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Legal Name: %s\n", result.LegalName)
    fmt.Printf("Status: %s\n", result.RegistrationStatus.Description)

    // Query a CPF through the third-party service cpfcnpj.com.br
    cpfResult, err := client.QueryCPFCNPJ(ctx, "529.982.247-25")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Name: %s\n", cpfResult.Name)
}
```

> 💡 **Cache**: Result caching is **disabled** by default. Enable it with
> `api.WithCacheTTL(duration)`. Only successful responses are cached, keyed by
> provider and document. The Serpro OAuth token is automatically reused across
> calls of the same `Client`.

### Generating documents

```go
package main

import (
    "fmt"

    "github.com/LeonardoSilva-10/brdocs"
)

func main() {
    cpf := brdocs.GenerateCPF()   // e.g.: "52998224806"
    cnpj := brdocs.GenerateCNPJ() // e.g.: "11222333000181"

    fmt.Println(brdocs.IsValidCPF(cpf))  // true
    fmt.Println(brdocs.IsValidCNPJ(cnpj)) // true
}
```

## 📋 Supported APIs

| API | Type | Document | Requires credentials | Free |
|-----|------|----------|----------------------|------|
| **Serpro/ConectaGov** | Official | CNPJ | ✅ OAuth2 (Consumer Key/Secret) | ❌ Paid |
| **Receita Federal** | Official | CPF | ❌ | ✅ Free (web scraping) |
| **cpfcnpj.com.br** | Third-party | CPF/CNPJ | ✅ API Key | ❌ Paid |
| **cpf-brasil.org** | Third-party | CPF | ✅ API Key | ❌ Paid (with trial) |

> ⚠️ **About CPF**: The Receita Federal does not offer a public official REST
> API to query CPF numbers. This package implements the query by parsing the
> public registration status page (web scraping).

## 🏗️ Project Structure

```
brdocs/
├── cpf/                  # CPF validation and generation
├── cnpj/                 # CNPJ validation and generation
├── api/                  # External API client
│   ├── client.go         # Client, options and cache
│   ├── serpro.go         # Serpro/ConectaGov integration (CNPJ)
│   ├── receitafederal.go # Receita Federal public query (CPF)
│   └── thirdparty.go     # cpfcnpj.com.br and cpf-brasil.org
├── brdocs.go             # Convenience re-exports
├── e2e/                  # End-to-end tests
└── .github/workflows/    # CI/CD
```

## 🧪 Testing

```bash
# All tests
go test ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# E2E tests (requires environment variables)
SERPRO_CONSUMER_KEY=xxx SERPRO_CONSUMER_SECRET=xxx go test ./e2e/... -run TestSerproIntegration
CPFCNPJ_API_KEY=xxx go test ./e2e/... -run TestCPFCNPJAPIIntegration

# Receita Federal public query (may break due to CAPTCHA/page changes)
RUN_RF_INTEGRATION=1 go test ./e2e/... -run TestReceitaFederalIntegration

# Benchmarks
go test -bench=. -benchmem ./...
```

## 📄 License

MIT License - see the [LICENSE](LICENSE) file for details.

---

**⚠️ Legal Notice**: This package is a development tool. Use of the official
Brazilian government APIs is subject to the Receita Federal and Serpro terms of
service. CPF data is protected by the LGPD (Brazilian General Data Protection
Law).