package models

// LoginRequest - dados para fazer login na Caixa
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse - resposta do login
type LoginResponse struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	SessionToken string `json:"session_token,omitempty"`
}

// SearchRequest - dados para buscar por CPF
type SearchRequest struct {
	CPF          string `json:"cpf"`
	SessionToken string `json:"session_token"`
}

// ClientData - dados extraídos do portal da Caixa
type ClientData struct {
	// Dados Pessoais do Proponente
	CPF               string `json:"cpf"`
	Nome              string `json:"nome"`
	Ocupacao          string `json:"ocupacao,omitempty"`
	Nacionalidade     string `json:"nacionalidade,omitempty"`
	TipoIdentificacao string `json:"tipo_identificacao,omitempty"`
	RG                string `json:"rg,omitempty"`
	
	// Dados de Contato
	TelefoneCelular string `json:"telefone_celular,omitempty"`
	
	// Dados do Contrato
	NumeroContrato        string `json:"numero_contrato"`
	AgendamentoAssinatura string `json:"agendamento_assinatura"`
	
	// Dados Bancários
	ContaDebitoCompleta string `json:"conta_debito_completa"`
	Agencia             string `json:"agencia"`
	ContaCorrente       string `json:"conta_corrente"`
	
	// Dados do Coobrigado
	CoobrigadoCPF  string `json:"coobrigado_cpf,omitempty"`
	CoobrigadoNome string `json:"coobrigado_nome,omitempty"`
	
	// Endereço Residencial do Proponente
	CEP            string `json:"cep,omitempty"`
	TipoLogradouro string `json:"tipo_logradouro,omitempty"`
	Logradouro     string `json:"logradouro,omitempty"`
	Numero         string `json:"numero,omitempty"`
	Bairro         string `json:"bairro,omitempty"`
	Municipio      string `json:"municipio,omitempty"`
	UF             string `json:"uf,omitempty"`
	Complemento    string `json:"complemento,omitempty"`
	
	// Dados do Imóvel
	EnderecoImovel string `json:"endereco_imovel,omitempty"`
	CEPImovel      string `json:"cep_imovel,omitempty"`
	
	// Dados Financeiros
	ValorCompraVenda string `json:"valor_compra_venda,omitempty"`
}

// SearchResponse - resposta da busca
type SearchResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *ClientData `json:"data,omitempty"`
}

// HealthResponse - resposta do health check
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// LoginAndSearchRequest - faz login e busca em uma única operação
type LoginAndSearchRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	CPF      string `json:"cpf"`
}