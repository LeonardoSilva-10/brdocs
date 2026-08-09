# Contributing to brdocs

Thanks for your interest in contributing!

## How to Contribute

### 1. Fork and Clone

```bash
git clone https://github.com/LeonardoSilva-10/brdocs.git
cd brdocs
```

### 2. Install Dependencies

```bash
make deps
make install-tools
```

### 3. Create a Branch

```bash
git checkout -b feature/my-feature
```

### 4. Develop and Test

```bash
# Run the tests
make test

# Run the linter
make lint

# Run the benchmarks
make bench
```

### 5. Commit and Push

```bash
git add .
git commit -m "feat: describe the feature"
git push origin feature/my-feature
```

### 6. Open a Pull Request

## Code Standards

- Follow the conventions of [Effective Go](https://go.dev/doc/effective_go)
- Document every exported function and type
- Keep test coverage above 80%
- Use `gofmt` for formatting

## Tests

### Unit Tests

```bash
go test ./...
```

### E2E Tests

```bash
# Local validation only
make e2e

# With Serpro integration (requires credentials)
SERPRO_CONSUMER_KEY=xxx SERPRO_CONSUMER_SECRET=xxx make e2e-serpro
```

### Benchmarks

```bash
make bench
```

## Reporting Bugs

Use the [GitHub Issues](https://github.com/LeonardoSilva-10/brdocs/issues) with:
- A clear description of the problem
- Steps to reproduce
- The Go version
- The error output

## License

By contributing, you agree that your contributions will be licensed under the MIT License.