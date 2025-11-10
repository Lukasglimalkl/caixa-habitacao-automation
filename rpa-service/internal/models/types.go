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
	SessionToken string `json:"session_token,omitempty"` // omitempty = só aparece se tiver valor
}

// SearchRequest - dados para buscar por CPF
type SearchRequest struct {
	CPF          string `json:"cpf"`
	SessionToken string `json:"session_token"`
}

// ClientData - dados extraídos do portal da Caixa
type ClientData struct {
	// Dados do Proponente (clicando no link)
	CPF                   string `json:"cpf"`
	Nome                  string `json:"nome"`
	NumeroContrato        string `json:"numero_contrato"`
	ContaDebitoCompleta   string `json:"conta_debito_completa"`
	Agencia               string `json:"agencia"`
	ContaCorrente         string `json:"conta_corrente"`
	AgendamentoAssinatura string `json:"agendamento_assinatura"`
	Endereco              string `json:"endereco,omitempty"`
	
	// Dados do Coobrigado (da tabela)
	CoobrigadoCPF  string `json:"coobrigado_cpf,omitempty"`
	CoobrigadoNome string `json:"coobrigado_nome,omitempty"`
}

// SearchResponse - resposta da busca
type SearchResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    *ClientData `json:"data,omitempty"` // ponteiro porque pode ser nil
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