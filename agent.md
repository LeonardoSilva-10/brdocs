# agent.md — Guia para agentes (e colaboradores) no brdocs

Este documento orienta qualquer pessoa — humana ou agente de IA — que vá tocar
neste repositório. Leia antes de editar código, testes ou documentação.

## 1. O que é o projeto

`brdocs` é uma **biblioteca Go** para validação, geração e consulta de CPF e
CNPJ brasileiros. Nada aqui é uma aplicação executável: é um módulo Go publicável
(`github.com/LeonardoSilva-10/brdocs`), então **a API pública importa** para quem
vai consumir.

- Validação local: sem dependências externas (só stdlib) nos pacotes `cpf` e `cnpj`.
- Consulta: `api` conversa com Serpro/ConectaGov, Receita Federal (scraping) e
  serviços de terceiros (cpfcnpj.com.br, cpf-brasil.org).
- Única dependência em `go.mod` é `testify` (usada **apenas em testes**).

## 2. Estrutura real (não invente pastas que não existem)

```
brdocs/
├── cpf/cpf.go              # Validação, formatação, máscara, UF e geração de CPF
├── cnpj/cnpj.go            # Validação, formatação, máscara, matriz/filial e geração de CNPJ
├── api/
│   ├── client.go           # Client, options, cache e sanitizeDocument
│   ├── serpro.go           # Serpro/ConectaGov (CNPJ) — OAuth2
│   ├── receitafederal.go   # Consulta pública da Receita Federal (CPF) — scraping
│   ├── thirdparty.go       # cpfcnpj.com.br e cpf-brasil.org
│   └── types.go            # Structs de resultado e helpers (IsActive, IsRegular)
├── brdocs.go               # Re-exports de conveniência (cpf/cnpj)
├── e2e/e2e_test.go         # Testes end-to-end (alguns exigem credenciais/env)
└── .github/workflows/ci.yml
```

> ⚠️ **Não** adicione pastas `internal/`, `mocks/` etc. sem necessidade real.
> O antigo `internal/httpclient` foi removido por ser código morto. Cada pacote
> dentro de `api/` é um arquivo `.go` no mesmo package `api`, não uma subpasta.

## 3. Convenções de código

- **Go idiomático**: siga `Effective Go` e `gofmt`. Rode `gofmt -w` antes de commitar.
- **Toda função/tipo/var exportado tem docstring** começando com o nome do símbolo.
- **Erros**: use `errors.New` para sentinelas (`ErrInvalidLength`, etc.) e
  `fmt.Errorf("pkg: ...", %w)` para enriquecer com contexto. Prefixe a mensagem
  com o pacote (`api:`, `serpro:`, `cpfcnpj:`).
- **Contexto**: toda chamada de API recebe `context.Context` como primeiro arg.
- **Testes**: use `testify` (`assert`/`require`). Prefira tabelas de teste.
- **Exemplos**: `ExampleXxx` funciona **como teste** — o bloco `// Output:` é
  verificado. Nunca use um CPF/CNPJ inválido num exemplo sem conferir o output.
- **Dependências**: não adicione dependências externas aos pacotes de validação
  (`cpf`/`cnpj`). Testes/API podem usar o que já existe (`testify`, stdlib).

## 4. Conhecimento de domínio (não erre nisso)

- **CPF**: 11 dígitos, 2 DV módulo 11 (pesos 10..2 e 11..2). Dígitos todos
  iguais (ex. `111.111.111-11`) são inválidos.
- **CNPJ**: 14 dígitos, 2 DV módulo 11 (pesos 5..2 e 6..2). Havendo dígitos
  repetidos, é inválido. Filial/matriz: posições 9–12 (`0001` = matriz).
- **Estado (UF) do CPF**: o **9º dígito** indica a UF de emissão. Um dígito pode
  mapear para várias UFs (retorna `[]string`). Ex.: dígito 2 → AC/AM/AP/PA/RO/RR.
- **Situação cadastral CNPJ (Receita Federal)** — tabela oficial:
  `1`=NULA, `2`=ATIVA, `3`=SUSPENSA, `4`=INAPTA, `8`=BAIXADA.
  Logo, `CNPJResult.IsActive()` SÓ considera o código `2` (ou descrição contendo
  "ATIVA"). **Não** trate `1` como ativo.
- **Situação CPF**: valores textuais (REGULAR / IRREGULAR / CANCELADA / NULA…).

## 5. API de consultas — regras ao mexer

- **Cache**: implementado em `api` via `WithCacheTTL` (desabilitado por padrão).
  Chaves são `provedor:documento` (Receita Federal inclui a data de nascimento).
  Só **respostas bem-sucedidas** são cacheadas. Ao adicionar provider, use
  `cacheGet`/`cacheSet` com prefixo próprio.
- **Token OAuth Serpro**: fica no `serproClient` **persistido no `Client`**
  (`c.serpro`) para ser reutilizado entre chamadas. Não recrie o cliente por chamada.
- **Sanitização**: use `sanitizeDocument` (remove não-dígitos) em vez de repetir
  `strings.ReplaceAll` em cada arquivo.
- **Resposta de terceiros (cpfcnpj)**: é tipada (`CPFCNPJResult`), não `interface{}`.
- **doRequest**: só decodifica 200 OK; erros de status retornam erro genérico.
  Não cache erros.

## 6. Áreas frágeis / conhecidas

- **Receita Federal (CPF)**: é **web scraping** da página pública. O site tem
  CAPTCHA e muda a estrutura HTML; os regexp de parsing podem quebrar a qualquer
  momento. É o motivo do teste E2E estar atrás de `RUN_RF_INTEGRATION=1`.
- **Validação do CPF via RF** depende de data de nascimento e exige o dígito de
  conferência interno. Não assuma que bater na URL funciona — é frágil.
- **`go test ./...` deve ficar verde SEMPRE**, inclusive com os exemplos
  (`ExampleState` já quebrou o pipeline por usar CPF inválido).

## 7. Comandos

```bash
go test ./...                # todos os testes (deve passar)
go test -race ./...          # com race detector
go test -cover ./...         # cobertura por pacote
go vet ./...                 # análise estática
gofmt -w $(find . -name '*.go')  # formatação
go test -bench=. -benchmem ./... # benchmarks

# E2E (exigem credenciais/ambiente; CI pula sem elas)
SERPRO_CONSUMER_KEY=x SERPRO_CONSUMER_SECRET=x go test ./e2e/... -run TestSerproIntegration
CPFCNPJ_API_KEY=x go test ./e2e/... -run TestCPFCNPJAPIIntegration
RUN_RF_INTEGRATION=1 go test ./e2e/... -run TestReceitaFederalIntegration
```

Se preferir, há um `Makefile` com `make test`, `make coverage`, `make lint`,
`make bench`, `make e2e`.

## 8. Checklist obrigatório para qualquer mudança

1. Compreenda o impacto na **API pública** (funções/structs exportados).
2. Rode `gofmt` e `go vet ./...`.
3. Rode `go test ./...` — **incluindo os `ExampleXxx`** — e veja tudo verde.
4. Se adicionar/alterar comportamento dos validadores, rode `go test -cover ./...`
   e evite regressão relevante de cobertura.
5. Atualize `README.md` (funcionalidades, exemplos, tabela de APIs, árvore) e
   `CHANGELOG.md` quando o comportamento visível mudar.
6. Mantenha o `agent.md` e o `README` **em sincronia com o código real**
   (pastas, nomes de função, covariáveis). Ex.: na raiz é `brdocs.IsValidCPF` e
   nos subpacotes é `cpf.Valid`/`cnpj.Valid` — não invente `*.IsValid`.
7. Não adicione credenciais/variáveis ao código — leia de environment nos testes
   E2E e pule (`t.Skip`) quando ausentes.

## 9. Privacidade e aviso legal

Dados de CPF são protegidos pela **LGPD**. Use CPFs/CNPJs de teste válidos e
não colete nem persista dados pessoais. As APIs oficiais têm termos de serviço
próprios (Serpro/Receita Federal). Documente qualquer limite ou exigência de
CAPTCHA/credenciais.

