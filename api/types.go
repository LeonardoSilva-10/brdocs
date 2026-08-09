package api

import (
	"strings"
	"time"
)

// RegistrationStatus represents the registration status.
type RegistrationStatus struct {
	Code        string `json:"codigo"`
	Description string `json:"descricao"`
	Date        string `json:"data"`
	Reason      string `json:"motivo"`
}

// Address represents a Brazilian address.
type Address struct {
	Street     string `json:"logradouro"`
	Number     string `json:"numero"`
	Complement string `json:"complemento"`
	PostalCode string `json:"cep"`
	District   string `json:"bairro"`
	City       string `json:"municipio"`
	State      string `json:"uf"`
}

// Phone represents a phone number.
type Phone struct {
	AreaCode string `json:"ddd"`
	Number   string `json:"numero"`
}

// CNAE represents a Brazilian economic activity code.
type CNAE struct {
	Code        string `json:"codigo"`
	Description string `json:"descricao"`
}

// LegalNature represents the legal nature of a company.
type LegalNature struct {
	Code        string `json:"codigo"`
	Description string `json:"descricao"`
}

// Partner represents a company partner/shareholder.
type Partner struct {
	CPF           string `json:"cpf"`
	Name          string `json:"nome"`
	Qualification string `json:"qualificacao"`
}

// CNPJResult holds the data returned from a CNPJ query.
type CNPJResult struct {
	Number             string             `json:"ni"`
	LegalName          string             `json:"nome_empresarial"`
	TradeName          string             `json:"nome_fantasia"`
	OpeningDate        string             `json:"data_abertura"`
	RegistrationStatus RegistrationStatus `json:"situacao_cadastral"`
	LegalNature        LegalNature        `json:"natureza_juridica"`
	PrimaryCNAE        CNAE               `json:"cnae_principal"`
	SecondaryCNAEs     []CNAE             `json:"cnae_secundarios"`
	Address            Address            `json:"endereco"`
	Phones             []Phone            `json:"telefones"`
	Email              string             `json:"correio_eletronico"`
	ShareCapital       float64            `json:"capital_social"`
	CompanySize        string             `json:"porte"`
	Partners           []Partner          `json:"socios"`
	SimplesOptedIn     bool               `json:"optante_simples"`
	SIMEIOptedIn       bool               `json:"optante_simei"`
	EstablishmentType  string             `json:"tipo_estabelecimento"`
	Agency             string             `json:"orgao"`
	FederativeEntity   string             `json:"ente_federativo"` //nolint:misspell // the json tag keeps the upstream field name "ente_federativo"
	SpecialStatus      string             `json:"situacao_especial"`
	SpecialStatusDate  string             `json:"data_situacao_especial"`
}

// CPFResult holds the data returned from a CPF query.
type CPFResult struct {
	CPF                string    `json:"cpf"`
	Name               string    `json:"nome"`
	CivilName          string    `json:"nome_civil"`
	SocialName         string    `json:"nome_social"`
	RegistrationStatus string    `json:"situacao_cadastral"`
	RegistrationDate   string    `json:"data_inscricao"`
	BirthDate          string    `json:"data_nascimento"`
	CheckDigit         string    `json:"digito_verificador"`
	DeathYear          string    `json:"ano_obito"`
	QueryTimestamp     time.Time `json:"consulta_datahora"`
}

// CPFCNPJResult holds the generic response returned by the cpfcnpj.com.br API.
type CPFCNPJResult struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Name      string `json:"nome"`
	LegalName string `json:"razao_social"`
	CPF       string `json:"cpf"`
	CNPJ      string `json:"cnpj"`
	BirthDate string `json:"data_nascimento"`
}

// IsActive returns true if the CNPJ has an active registration status.
//
// The Receita Federal "situação cadastral" code for CNPJ is 2 (ATIVA). Codes
// 1 (NULA), 3 (SUSPENSA), 4 (INAPTA) and 8 (BAIXADA) are not active. The
// description is also matched as a fallback for providers that return
// descriptive values instead of codes.
func (r *CNPJResult) IsActive() bool {
	return r.RegistrationStatus.Code == "2" ||
		strings.Contains(strings.ToUpper(r.RegistrationStatus.Description), "ATIVA")
}

// IsRegular returns true if the CPF is regular.
func (r *CPFResult) IsRegular() bool {
	return r.RegistrationStatus == "0" || r.RegistrationStatus == "Regular"
}
