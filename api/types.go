package api

import "time"

// SituacaoCadastral represents the registration status.
type SituacaoCadastral struct {
	Codigo     string `json:"codigo"`
	Descricao  string `json:"descricao"`
	Data       string `json:"data"`
	Motivo     string `json:"motivo"`
}

// Endereco represents a Brazilian address.
type Endereco struct {
	Logradouro  string `json:"logradouro"`
	Numero      string `json:"numero"`
	Complemento string `json:"complemento"`
	CEP         string `json:"cep"`
	Bairro      string `json:"bairro"`
	Municipio   string `json:"municipio"`
	UF          string `json:"uf"`
}

// Telefone represents a phone number.
type Telefone struct {
	DDD    string `json:"ddd"`
	Numero string `json:"numero"`
}

// CNAE represents a Brazilian economic activity code.
type CNAE struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
}

// NaturezaJuridica represents the legal nature of a company.
type NaturezaJuridica struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
}

// Socio represents a company partner/shareholder.
type Socio struct {
	CPF           string `json:"cpf"`
	Nome          string `json:"nome"`
	Qualificacao  string `json:"qualificacao"`
}

// CNPJResult holds the data returned from a CNPJ query.
type CNPJResult struct {
	NI                  string            `json:"ni"`
	NomeEmpresarial     string            `json:"nome_empresarial"`
	NomeFantasia        string            `json:"nome_fantasia"`
	DataAbertura        string            `json:"data_abertura"`
	SituacaoCadastral   SituacaoCadastral `json:"situacao_cadastral"`
	NaturezaJuridica    NaturezaJuridica  `json:"natureza_juridica"`
	CNAEPrincipal       CNAE              `json:"cnae_principal"`
	CNAESecundarios     []CNAE            `json:"cnae_secundarios"`
	Endereco            Endereco          `json:"endereco"`
	Telefones           []Telefone        `json:"telefones"`
	CorreioEletronico   string            `json:"correio_eletronico"`
	CapitalSocial       float64           `json:"capital_social"`
	Porte               string            `json:"porte"`
	Socios              []Socio           `json:"socios"`
	OptanteSimples      bool              `json:"optante_simples"`
	OptanteSIMEI        bool              `json:"optante_simei"`
	TipoEstabelecimento string            `json:"tipo_estabelecimento"`
	Orgao               string            `json:"orgao"`
	EnteFederativo      string            `json:"ente_federativo"`
	SituacaoEspecial    string            `json:"situacao_especial"`
	DataSituacaoEspecial string           `json:"data_situacao_especial"`
}

// CPFResult holds the data returned from a CPF query.
type CPFResult struct {
	CPF                string    `json:"cpf"`
	Nome               string    `json:"nome"`
	NomeCivil          string    `json:"nome_civil"`
	NomeSocial         string    `json:"nome_social"`
	SituacaoCadastral  string    `json:"situacao_cadastral"`
	DataInscricao      string    `json:"data_inscricao"`
	DataNascimento     string    `json:"data_nascimento"`
	DigitoVerificador  string    `json:"digito_verificador"`
	AnoObito           string    `json:"ano_obito"`
	ConsultaDataHora  time.Time `json:"consulta_datahora"`
}

// IsActive returns true if the CNPJ is active (situação cadastral = 1).
func (r *CNPJResult) IsActive() bool {
	return r.SituacaoCadastral.Codigo == "1" || r.SituacaoCadastral.Codigo == "2"
}

// IsRegular returns true if the CPF is regular.
func (r *CPFResult) IsRegular() bool {
	return r.SituacaoCadastral == "0" || r.SituacaoCadastral == "Regular"
}
