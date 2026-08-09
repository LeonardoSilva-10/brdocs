package api

import (
	"context"
	"fmt"
	"net/http"
)

// thirdPartyClient handles third-party API integrations.
type thirdPartyClient struct {
	parent *Client
}

func newThirdPartyClient(parent *Client) *thirdPartyClient {
	return &thirdPartyClient{parent: parent}
}

// QueryCPFCNPJ queries CPF or CNPJ from cpfcnpj.com.br.
func (tp *thirdPartyClient) QueryCPFCNPJ(ctx context.Context, document string) (*CPFCNPJResult, error) {
	if tp.parent.cpfCnpjAPIKey == "" {
		return nil, fmt.Errorf("cpfcnpj: API key is required")
	}

	document = sanitizeDocument(document)
	cacheKey := "cpfcnpj:" + document
	if cached := tp.parent.cacheGet(cacheKey); cached != nil {
		return cached.(*CPFCNPJResult), nil
	}

	endpoint := fmt.Sprintf("%s/%s/%s", tp.parent.baseURLs["cpfcnpj"], document, tp.parent.cpfCnpjAPIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("cpfcnpj: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	var result CPFCNPJResult
	if err := tp.parent.doRequest(ctx, req, &result); err != nil {
		return nil, fmt.Errorf("cpfcnpj: %w", err)
	}

	if result.Status != "" && result.Status != "OK" {
		return nil, fmt.Errorf("cpfcnpj: query failed: %s", result.Message)
	}

	tp.parent.cacheSet(cacheKey, &result)
	return &result, nil
}

// QueryCPFBrasil queries CPF from cpf-brasil.org.
func (tp *thirdPartyClient) QueryCPFBrasil(ctx context.Context, cpf string) (*CPFResult, error) {
	if tp.parent.cpfBrasilAPIKey == "" {
		return nil, fmt.Errorf("cpfbrasil: API key is required")
	}

	cpf = sanitizeDocument(cpf)
	cacheKey := "cpfbrasil:" + cpf
	if cached := tp.parent.cacheGet(cacheKey); cached != nil {
		return cached.(*CPFResult), nil
	}

	endpoint := fmt.Sprintf("%s/cpf/%s", tp.parent.baseURLs["cpfbrasil"], cpf)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("cpfbrasil: failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", tp.parent.cpfBrasilAPIKey)
	req.Header.Set("Accept", "application/json")

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			CPF  string `json:"CPF"`
			Name string `json:"NOME"`
			Sexo string `json:"SEXO"`
			Nasc string `json:"NASC"`
		} `json:"data"`
	}

	if err := tp.parent.doRequest(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("cpfbrasil: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("cpfbrasil: query failed")
	}

	result := &CPFResult{
		CPF:       response.Data.CPF,
		Name:      response.Data.Name,
		BirthDate: response.Data.Nasc,
	}
	tp.parent.cacheSet(cacheKey, result)
	return result, nil
}

// QueryCPFCNPJ is a convenience method on Client.
func (c *Client) QueryCPFCNPJ(ctx context.Context, document string) (*CPFCNPJResult, error) {
	tp := newThirdPartyClient(c)
	return tp.QueryCPFCNPJ(ctx, document)
}

// QueryCPFBrasil is a convenience method on Client.
func (c *Client) QueryCPFBrasil(ctx context.Context, cpf string) (*CPFResult, error) {
	tp := newThirdPartyClient(c)
	return tp.QueryCPFBrasil(ctx, cpf)
}
