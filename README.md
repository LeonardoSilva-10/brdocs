# brdocs

[![Go Reference](https://pkg.go.dev/badge/github.com/LeonardoSilva-10/brdocs.svg)](https://pkg.go.dev/github.com/LeonardoSilva-10/brdocs)
[![CI](https://github.com/LeonardoSilva-10/brdocs/actions/workflows/ci.yml/badge.svg)](https://github.com/LeonardoSilva-10/brdocs/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/LeonardoSilva-10/brdocs)](https://goreportcard.com/report/github.com/LeonardoSilva-10/brdocs)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> Pacote Go para validação e consulta de CPF e CNPJ brasileiros.

## ✨ Funcionalidades

- ✅ **Validação matemática** de CPF e CNPJ (algoritmo oficial módulo 11)
- 🔍 **Consulta via APIs** do governo brasileiro e serviços de terceiros
- 🛡️ **Sanitização** automática (remove pontos, traços, barras)
- 🎲 **Geração** de CPF e CNPJ válidos (`crypto/rand`)
- 📦 **Sem dependências externas** para validação local
- 🧪 **Testes unitários e E2E** com alta cobertura nos validadores (~90%)
- ⚡ **Cache embutido** para consultas repetidas (`api.WithCacheTTL`)
- 🌐 **Contexto nativo** (`context.Context`) em todas as chamadas de API

## 📦 Instalação

```bash
go get github.com/LeonardoSilva-10/brdocs
```

## 🚀 Uso Rápido

### Validação Local (sem internet)

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
    fmt.Println(cpf.Valid("111.111.111-11")) // false (digitos repetidos)
    fmt.Println(cpf.Valid("123.456.789-09")) // false (digitos verificadores invalidos)

    // CNPJ
    fmt.Println(cnpj.Valid("11.222.333/0001-81")) // true
    fmt.Println(cnpj.Valid("00.000.000/0000-00")) // false
}
```

### Consulta via API

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
        api.WithCacheTTL(5*time.Minute), // cache resultados de consultas por 5 minutos
    )

    // Consulta CNPJ na Receita Federal (Serpro)
    result, err := client.ConsultaCNPJ(ctx, "34238864000168")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Razão Social: %s\n", result.NomeEmpresarial)
    fmt.Printf("Situação: %s\n", result.SituacaoCadastral.Descricao)

    // Consulta CPF no serviço de terceiros cpfcnpj.com.br
    cpfResult, err := client.ConsultaCPFCNPJ(ctx, "529.982.247-25")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Nome: %s\n", cpfResult.Nome)
}
```

> 💡 **Cache**: Por padrão o cache de consultas está **desabilitado**. Habilite-o
> com `api.WithCacheTTL(duration)`. Apenas respostas bem-sucedidas são
> armazenadas, separadas por provedor e documento. O token OAuth do Serpro já é
> reutilizado automaticamente entre chamadas do mesmo `Client`.

### Geração de documentos

```go
package main

import (
    "fmt"

    "github.com/LeonardoSilva-10/brdocs"
)

func main() {
    cpf := brdocs.GenerateCPF()   // ex.: "52998224806"
    cnpj := brdocs.GenerateCNPJ() // ex.: "11222333000181"

    fmt.Println(brdocs.IsValidCPF(cpf))  // true
    fmt.Println(brdocs.IsValidCNPJ(cnpj)) // true
}
```

## 📋 APIs Suportadas

| API | Tipo | Documento | Requer Credenciais | Gratuita |
|-----|------|-----------|-------------------|----------|
| **Serpro/ConectaGov** | Oficial | CNPJ | ✅ OAuth2 (Consumer Key/Secret) | ❌ Paga |
| **Receita Federal** | Oficial | CPF | ❌ | ✅ Gratuita (web scraping) |
| **cpfcnpj.com.br** | Terceiro | CPF/CNPJ | ✅ API Key | ❌ Paga |
| **cpf-brasil.org** | Terceiro | CPF | ✅ API Key | ❌ Paga (com trial) |

> ⚠️ **Nota sobre CPF**: A Receita Federal não oferece uma API REST pública oficial para consulta de CPF. O pacote implementa consulta via parsing da página pública de situação cadastral.

## 🏗️ Estrutura do Projeto

```
brdocs/
├── cpf/                  # Validação e geração de CPF
├── cnpj/                 # Validação e geração de CNPJ
├── api/                  # Cliente de APIs externas
│   ├── client.go         # Cliente, opções e cache
│   ├── serpro.go         # Integração Serpro/ConectaGov (CNPJ)
│   ├── receitafederal.go # Consulta pública Receita Federal (CPF)
│   └── thirdparty.go     # cpfcnpj.com.br e cpf-brasil.org
├── brdocs.go             # Re-exports de conveniência
├── e2e/                  # Testes end-to-end
└── .github/workflows/    # CI/CD
```

## 🧪 Testes

```bash
# Todos os testes
go test ./...

# Com cobertura
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Testes E2E (requer variáveis de ambiente)
SERPRO_CONSUMER_KEY=xxx SERPRO_CONSUMER_SECRET=xxx go test ./e2e/... -run TestSerproIntegration
CPFCNPJ_API_KEY=xxx go test ./e2e/... -run TestCPFCNPJAPIIntegration

# Consulta pública da Receita Federal (pode quebrar por CAPTCHA/mudanças na página)
RUN_RF_INTEGRATION=1 go test ./e2e/... -run TestReceitaFederalIntegration

# Benchmarks
go test -bench=. -benchmem ./...
```

## 📄 Licença

MIT License - veja [LICENSE](LICENSE) para detalhes.

---

**⚠️ Aviso Legal**: Este pacote é uma ferramenta de desenvolvimento. O uso de APIs oficiais do governo brasileiro está sujeito aos termos de serviço da Receita Federal e Serpro. Dados de CPF são protegidos pela LGPD.