package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// Handler - gerencia as requisições HTTP
type Handler struct {
	bot *automation.CaixaBot
}

// NewHandler - cria um novo handler
func NewHandler() *Handler {
	return &Handler{
		bot: automation.NewCaixaBot(true), // headless = true
	}
}

// LoginAndSearch - endpoint para login e busca
func (h *Handler) LoginAndSearch(w http.ResponseWriter, r *http.Request) {
	var req models.LoginAndSearchRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	logger.Info("📥 Nova requisição recebida")
	logger.Info("👤 Usuário: " + req.Username)
	logger.Info("🔍 CPF: " + req.CPF)
	
	// Executa automação
	response, err := h.bot.LoginAndSearch(req.Username, req.Password, req.CPF)
	
	w.Header().Set("Content-Type", "application/json")
	
	if err != nil {
		logger.Error("❌ Erro na automação: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	logger.Info("✅ Requisição processada com sucesso!")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Health - endpoint de health check
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:  "healthy",
		Service: "RPA Service - Caixa Automation",
		Version: "2.0.0",
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}