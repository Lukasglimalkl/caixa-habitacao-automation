package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// Handler - estrutura que vai conter os handlers HTTP
type Handler struct {
	bot *automation.CaixaBot
}

// NewHandler - cria um novo handler
func NewHandler() *Handler {
	return &Handler{
		bot: automation.NewCaixaBot(),
	}
}

// HealthCheck - endpoint de health check
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	logger.Info("Health check requisitado")

	response := models.HealthResponse{
		Status:  "UP",
		Service: "RPA Service - Caixa Automation",
		Version: "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Login - endpoint para fazer login na Caixa
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	logger.Info("Requisição de login recebida")

	// Apenas aceita POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o JSON do body
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar request: " + err.Error())
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Valida campos obrigatórios
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username e password são obrigatórios", http.StatusBadRequest)
		return
	}

	// Chama o bot para fazer login
	response, err := h.bot.Login(req)
	if err != nil {
		logger.Error("Erro no login: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Retorna sucesso
	w.Header().Set("Content-Type", "application/json")
	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
	json.NewEncoder(w).Encode(response)
}

// SearchByCPF - endpoint para buscar por CPF
func (h *Handler) SearchByCPF(w http.ResponseWriter, r *http.Request) {
	logger.Info("Requisição de busca por CPF recebida")

	// Apenas aceita POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	// Decodifica o JSON do body
	var req models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar request: " + err.Error())
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Valida campos obrigatórios
	if req.CPF == "" {
		http.Error(w, "CPF é obrigatório", http.StatusBadRequest)
		return
	}

	if req.SessionToken == "" {
		http.Error(w, "Session token é obrigatório", http.StatusBadRequest)
		return
	}

	// Chama o bot para buscar
	response, err := h.bot.SearchByCPF(req)
	if err != nil {
		logger.Error("Erro na busca: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Retorna sucesso
	w.Header().Set("Content-Type", "application/json")
	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(response)
	
}

// LoginAndSearch - endpoint que faz login e busca em uma única operação
func (h *Handler) LoginAndSearch(w http.ResponseWriter, r *http.Request) {
	logger.Info("Requisição de login + busca recebida")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginAndSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Erro ao decodificar request: " + err.Error())
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Valida campos
	if req.Username == "" || req.Password == "" || req.CPF == "" {
		http.Error(w, "Username, password e CPF são obrigatórios", http.StatusBadRequest)
		return
	}

	// Chama o bot
	response, err := h.bot.LoginAndSearch(req)
	if err != nil {
		logger.Error("Erro no processo: " + err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if response.Success {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(response)
}

