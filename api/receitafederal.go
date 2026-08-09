package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// receitaFederalClient queries the public Receita Federal website for CPF data.
// Note: This is NOT an official REST API. It parses the public consultation page.
type receitaFederalClient struct {
	parent *Client
}

func newReceitaFederalClient(parent *Client) *receitaFederalClient {
	return &receitaFederalClient{parent: parent}
}

// QueryCPF queries CPF data from the public Receita Federal website.
// Requires the CPF number and birth date (DD/MM/YYYY format).
func (rf *receitaFederalClient) QueryCPF(ctx context.Context, cpf, birthDate string) (*CPFResult, error) {
	// Sanitize CPF
	cpf = sanitizeDocument(cpf)

	if len(cpf) != 11 {
		return nil, fmt.Errorf("receitafederal: invalid CPF length")
	}

	cacheKey := "receitafederal:" + cpf + ":" + birthDate
	if cached := rf.parent.cacheGet(cacheKey); cached != nil {
		return cached.(*CPFResult), nil
	}

	formData := url.Values{}
	formData.Set("txtCPF", cpf)
	formData.Set("txtDataNascimento", birthDate)
	formData.Set("Enviar", "Consultar")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		rf.parent.baseURLs["receitafederal"]+"/Servicos/CPF/ConsultaSituacao/ConsultaPublica.asp",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return nil, fmt.Errorf("receitafederal: failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; brdocs/1.0)")

	resp, err := rf.parent.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("receitafederal: request failed: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("receitafederal: unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("receitafederal: failed to read response: %w", err)
	}

	result, err := parseCPFResponse(string(body), cpf)
	if err != nil {
		return nil, err
	}

	rf.parent.cacheSet(cacheKey, result)
	return result, nil
}

// parseCPFResponse extracts CPF data from the HTML response.
func parseCPFResponse(html, cpf string) (*CPFResult, error) {
	result := &CPFResult{CPF: cpf}

	// Extract name
	nameRe := regexp.MustCompile(`<span[^>]*class="[^"]*nome[^"]*"[^>]*>([^<]+)</span>`)
	if matches := nameRe.FindStringSubmatch(html); len(matches) > 1 {
		result.Name = strings.TrimSpace(matches[1])
	}

	// Extract registration status
	statusRe := regexp.MustCompile(`Situação Cadastral[^>]*>\s*<[^>]*>([^<]+)</`)
	if matches := statusRe.FindStringSubmatch(html); len(matches) > 1 {
		result.RegistrationStatus = strings.TrimSpace(matches[1])
	}

	// Alternative patterns for different page versions
	if result.RegistrationStatus == "" {
		statusRe2 := regexp.MustCompile(`SITUAÇÃO CADASTRAL[^>]*>\s*<[^>]*>([^<]+)</`)
		if matches := statusRe2.FindStringSubmatch(html); len(matches) > 1 {
			result.RegistrationStatus = strings.TrimSpace(matches[1])
		}
	}

	// Extract registration date
	dataRe := regexp.MustCompile(`Data de Inscrição[^>]*>\s*<[^>]*>([^<]+)</`)
	if matches := dataRe.FindStringSubmatch(html); len(matches) > 1 {
		result.RegistrationDate = strings.TrimSpace(matches[1])
	}

	// Extract check digit
	dvRe := regexp.MustCompile(`Dígito Verificador[^>]*>\s*<[^>]*>([^<]+)</`)
	if matches := dvRe.FindStringSubmatch(html); len(matches) > 1 {
		result.CheckDigit = strings.TrimSpace(matches[1])
	}

	// If we couldn't extract any meaningful data, the CPF might be invalid or the page changed
	if result.Name == "" && result.RegistrationStatus == "" {
		// Check for error messages
		errorRe := regexp.MustCompile(`erro|inválido|não encontrado|não existe`)
		if errorRe.MatchString(strings.ToLower(html)) {
			return nil, fmt.Errorf("receitafederal: CPF not found or invalid")
		}
		return nil, fmt.Errorf("receitafederal: unable to parse response, page structure may have changed")
	}

	result.QueryTimestamp = time.Now()
	return result, nil
}

// QueryCPF is a convenience method on Client.
func (c *Client) QueryCPF(ctx context.Context, cpf, birthDate string) (*CPFResult, error) {
	rf := newReceitaFederalClient(c)
	return rf.QueryCPF(ctx, cpf, birthDate)
}
