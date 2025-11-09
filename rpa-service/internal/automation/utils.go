package automation

import (
	"fmt"
	"strings"

	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// separarContaDebito - separa agência e conta corrente
// Formato esperado: 0347-3701-000573937131-3
// 0347 = primeiros 4 dígitos
// 3701 = agência (próximos 4 dígitos após o primeiro traço)
// 000573937131-3 = conta corrente (restante)
func separarContaDebito(contaCompleta string) (agencia, contaCorrente string) {
	logger.Info(fmt.Sprintf("🔧 Separando conta: %s", contaCompleta))

	// Remove espaços
	contaCompleta = strings.TrimSpace(contaCompleta)

	// Divide por traço
	partes := strings.Split(contaCompleta, "-")

	if len(partes) >= 3 {
		// Agência é a segunda parte (índice 1)
		agencia = partes[1]
		
		// Conta corrente é a terceira parte em diante, juntando com traço
		contaCorrente = strings.Join(partes[2:], "-")
		
		logger.Info(fmt.Sprintf("✓ Agência: %s | Conta: %s", agencia, contaCorrente))
	} else {
		logger.Error("❌ Formato de conta inválido!")
		agencia = ""
		contaCorrente = contaCompleta
	}

	return agencia, contaCorrente
}