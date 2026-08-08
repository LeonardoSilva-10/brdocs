// Package api provides HTTP clients for querying Brazilian CPF and CNPJ data
// from external APIs.
//
// Supported providers:
//   - Serpro/ConectaGov (official, requires OAuth2 credentials)
//   - Receita Federal (public CPF query, no credentials required)
//   - cpfcnpj.com.br (third-party, requires API key)
//   - cpf-brasil.org (third-party, requires API key)
//
// Example:
//
//	client := api.NewClient(
//	    api.WithSerproCredentials("key", "secret"),
//	)
//	result, err := client.ConsultaCNPJ(ctx, "34238864000168")
package api
