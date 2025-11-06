package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/handlers"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
	"github.com/rs/cors"
)

func main() {
	// Inicializa o logger
	logger.Init()
	logger.Info("🚀 Iniciando RPA Service - Caixa Automation")

	// Cria o handler
	handler := handlers.NewHandler()

	// Configura as rotas
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Rota principal - Login + Busca
	router.HandleFunc("/api/login-and-search", handler.LoginAndSearch).Methods("POST")

	// Configura CORS (permite requisições do backend)
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Em produção, coloque apenas o domínio do backend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Aplica CORS no router
	httpHandler := corsHandler.Handler(router)

	// Porta do servidor
	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	// Inicia o servidor em uma goroutine
	go func() {
		logger.Info(fmt.Sprintf("🌐 Servidor rodando em http://localhost:%s", port))
		logger.Info("📋 Endpoints disponíveis:")
		logger.Info("   GET  /health                - Health check")
		logger.Info("   POST /api/login-and-search  - Login + Busca CPF (COMPLETO)")

		if err := http.ListenAndServe(addr, httpHandler); err != nil {
			logger.Error(fmt.Sprintf("Erro ao iniciar servidor: %v", err))
			os.Exit(1)
		}
	}()

	// Graceful shutdown - espera por CTRL+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 Encerrando servidor...")
	logger.Info("✅ Servidor encerrado com sucesso")
}

// getEnv - pega variável de ambiente ou retorna valor padrão
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}