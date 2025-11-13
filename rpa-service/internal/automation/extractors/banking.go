package extractors

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// CaixaBankingExtractor - implementação para extração bancária
type CaixaBankingExtractor struct{}

// NewBankingExtractor - cria novo extrator bancário
func NewBankingExtractor() *CaixaBankingExtractor {
	return &CaixaBankingExtractor{}
}

// ExtractBankingData - extrai dados bancários (conta de débito)
func (e *CaixaBankingExtractor) ExtractBankingData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("💳 Extraindo dados bancários...")
	
	// Faz scroll até a tabela de conta
	ScrollToTable(ctx, "Dados da Conta - Débito")
	
	// Extrai conta de débito
	return e.extractContaDebito(ctx, iframeNode, clientData)
}

// extractContaDebito - extrai a conta de débito completa
func (e *CaixaBankingExtractor) extractContaDebito(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🔍 Extraindo Conta de Débito...")
	
	// XPath específico para encontrar a conta de débito
	xpath := `//tr[@class='linha_azul'][.//label[contains(., 'Conta de Débito:')]]/td[@class='alinha_esquerda fonte_laranja']`
	
	var contaDebito string
	err := chromedp.Text(xpath, &contaDebito, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro com XPath principal: %v", err))
		logger.Info("🔍 Tentando XPath alternativo...")
		
		// XPath alternativo
		xpathAlt := `//td[@class='alinha_esquerda fonte_laranja' and contains(text(), '-') and contains(text(), '0347')]`
		err = chromedp.Text(xpathAlt, &contaDebito, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
		
		if err != nil {
			logger.Error(fmt.Sprintf("❌ XPath alternativo também falhou: %v", err))
			return fmt.Errorf("conta de débito não encontrada")
		}
	}
	
	clientData.ContaDebitoCompleta = strings.TrimSpace(contaDebito)
	logger.Info(fmt.Sprintf("✓ Conta completa: %s", clientData.ContaDebitoCompleta))
	
	// Separa agência e conta
	if clientData.ContaDebitoCompleta != "" {
		agencia, conta := e.separarContaDebito(clientData.ContaDebitoCompleta)
		clientData.Agencia = agencia
		clientData.ContaCorrente = conta
		logger.Info(fmt.Sprintf("✓ Agência: %s | Conta: %s", agencia, conta))
	} else {
		logger.Error("❌ Conta de débito está vazia!")
	}
	
	return nil
}

// separarContaDebito - separa agência e conta corrente
// Formato esperado: "0347-3701-000573937131-3" ou similar
func (e *CaixaBankingExtractor) separarContaDebito(contaCompleta string) (agencia, conta string) {
	logger.Info(fmt.Sprintf("🔧 Separando conta: %s", contaCompleta))
	
	// Remove espaços
	contaCompleta = strings.TrimSpace(contaCompleta)
	
	// Split por hífen
	partes := strings.Split(contaCompleta, "-")
	
	if len(partes) >= 3 {
		// Formato: 0347-3701-000573937131-3
		agencia = partes[1] // 3701
		// Conta = resto (pode ter mais de um hífen)
		conta = strings.Join(partes[2:], "-") // 000573937131-3
	} else if len(partes) == 2 {
		// Formato alternativo: 3701-000573937131-3
		agencia = partes[0]
		conta = partes[1]
	} else {
		// Formato desconhecido, tenta pegar tudo
		agencia = ""
		conta = contaCompleta
	}
	
	return agencia, conta
}