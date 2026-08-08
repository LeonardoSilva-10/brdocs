# Contributing to brdocs

Obrigado pelo interesse em contribuir!

## Como Contribuir

### 1. Fork e Clone

```bash
git clone https://github.com/LeonardoSilva-10/brdocs.git
cd brdocs
```

### 2. Instale as Dependências

```bash
make deps
make install-tools
```

### 3. Crie uma Branch

```bash
git checkout -b feature/minha-feature
```

### 4. Desenvolva e Teste

```bash
# Execute os testes
make test

# Execute o linter
make lint

# Execute os benchmarks
make bench
```

### 5. Commit e Push

```bash
git add .
git commit -m "feat: descrição da feature"
git push origin feature/minha-feature
```

### 6. Abra um Pull Request

## Padrões de Código

- Siga as convenções do [Effective Go](https://go.dev/doc/effective_go)
- Documente todas as funções e tipos exportados
- Mantenha a cobertura de testes acima de 80%
- Use `gofmt` para formatação

## Testes

### Testes Unitários

```bash
go test ./...
```

### Testes E2E

```bash
# Apenas validação local
make e2e

# Com integração Serpro (requer credenciais)
SERPRO_CONSUMER_KEY=xxx SERPRO_CONSUMER_SECRET=xxx make e2e-serpro
```

### Benchmarks

```bash
make bench
```

## Reportando Bugs

Use as [Issues do GitHub](https://github.com/LeonardoSilva-10/brdocs/issues) com:
- Descrição clara do problema
- Passos para reproduzir
- Versão do Go
- Saída do erro

## Licença

Ao contribuir, você concorda que suas contribuições serão licenciadas sob a MIT License.
