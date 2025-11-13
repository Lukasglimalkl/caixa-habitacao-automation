package main

import (
	"fmt"

	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

func main() {
	// IMPORTANTE: Inicializa o logger
	logger.Init()
	
	fmt.Println("🧪 Testando automação diretamente...")
	
	bot := automation.NewCaixaBot(false) // headless = false
	
	response, err := bot.LoginAndSearch(
		"marcella@iprimeconsultoria.com.br",
		"Bolinho2020",
		"53304153810",
	)
	
	if err != nil {
		fmt.Printf("❌ Erro: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Sucesso!\n")
	fmt.Printf("📋 Nome: %s\n", response.Data.Nome)
	fmt.Printf("📋 CPF: %s\n", response.Data.CPF)
}