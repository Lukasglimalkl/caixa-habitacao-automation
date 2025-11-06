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
// ClientData - dados extraídos do portal da Caixa
type ClientData struct {
	CPF                      string `json:"cpf"`
	Nome                     string `json:"nome"`
	Endereco                 string `json:"endereco"`
	AgendamentoAssinatura    string `json:"agendamento_assinatura"`    // 👈 NOVO
	NumeroContrato           string `json:"numero_contrato,omitempty"` // 👈 NOVO
	// Adicione mais campos conforme você descobre no scraping
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