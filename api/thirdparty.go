package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// thirdPartyClient handles third-party API integrations.
type thirdPartyClient struct {
	parent *Client
}

func newThirdPartyClient(parent *Client) *thirdPartyClient {
	return &thirdPartyClient{parent: parent}
}

// ConsultaCPFCPFCNPJ queries CPF or CNPJ from cpfcnpj.com.br.
func (tp *thirdPartyClient) ConsultaCPFCPFCNPJ(ctx context.Context, document string) (interface{}, error) {
	if tp.parent.cpfCnpjAPIKey == "" {
		return nil, fmt.Errorf("cpfcnpj: API key is required")
	}

	document = strings.ReplaceAll(document, ".", "")
	document = strings.ReplaceAll(document, "/", "")
	document = strings.ReplaceAll(document, "-", "")

	endpoint := fmt.Sprintf("%s/%s/%s", tp.parent.baseURLs["cpfcnpj"], document, tp.parent.cpfCnpjAPIKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("cpfcnpj: failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	var result map[string]interface{}
	if err := tp.parent.doRequest(ctx, req, &result); err != nil {
		return nil, fmt.Errorf("cpfcnpj: %w", err)
	}

	return result, nil
}

// ConsultaCPFBrasil queries CPF from cpf-brasil.org.
func (tp *thirdPartyClient) ConsultaCPFBrasil(ctx context.Context, cpf string) (*CPFResult, error) {
	if tp.parent.cpfBrasilAPIKey == "" {
		return nil, fmt.Errorf("cpfbrasil: API key is required")
	}

	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")

	endpoint := fmt.Sprintf("%s/cpf/%s", tp.parent.baseURLs["cpfbrasil"], cpf)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("cpfbrasil: failed to create request: %w", err)
	}

	req.Header.Set("X-API-Key", tp.parent.cpfBrasilAPIKey)
	req.Header.Set("Accept", "application/json")

	var response struct {
		Success bool `json:"success"`
		Data    struct {
			CPF    string `json:"CPF"`
			Nome   string `json:"NOME"`
			Sexo   string `json:"SEXO"`
			Nasc   string `json:"NASC"`
		} `json:"data"`
	}

	if err := tp.parent.doRequest(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("cpfbrasil: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("cpfbrasil: query failed")
	}

	return &CPFResult{
		CPF:            response.Data.CPF,
		Nome:           response.Data.Nome,
		DataNascimento: response.Data.Nasc,
	}, nil
}

// ConsultaCPFCPFCNPJ is a convenience method on Client.
func (c *Client) ConsultaCPFCPFCNPJ(ctx context.Context, document string) (interface{}, error) {
	tp := newThirdPartyClient(c)
	return tp.ConsultaCPFCPFCNPJ(ctx, document)
}

// ConsultaCPFBrasil is a convenience method on Client.
func (c *Client) ConsultaCPFBrasil(ctx context.Context, cpf string) (*CPFResult, error) {
	tp := newThirdPartyClient(c)
	return tp.ConsultaCPFBrasil(ctx, cpf)
}
