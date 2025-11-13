package extractors

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// CaixaContactExtractor - implementação para dados de contato
type CaixaContactExtractor struct{}

// NewContactExtractor - cria novo extrator de contato
func NewContactExtractor() *CaixaContactExtractor {
	return &CaixaContactExtractor{}
}

// ExtractContactData - extrai dados de contato (telefone)
func (e *CaixaContactExtractor) ExtractContactData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("📱 Extraindo dados de contato...")
	
	telefone := e.tryExtractTelefone(ctx, iframeNode)
	clientData.TelefoneCelular = telefone
	
	return nil
}

// tryExtractTelefone - tenta extrair telefone com múltiplos fallbacks
func (e *CaixaContactExtractor) tryExtractTelefone(ctx context.Context, iframeNode *cdp.Node) string {
	// Prioridade 1: Telefone Celular
	telefoneCelular, err := ExtractFieldFromTable(ctx, iframeNode, "Telefone Celular:")
	if err == nil && telefoneCelular != "" {
		logger.Info(fmt.Sprintf("✓ Telefone Celular: %s", telefoneCelular))
		return telefoneCelular
	}
	
	// Prioridade 2: Telefone Residencial
	logger.Info("⚠️ Telefone Celular vazio, tentando Residencial...")
	telefoneResidencial, err := ExtractFieldFromTable(ctx, iframeNode, "Telefone Residencial:")
	if err == nil && telefoneResidencial != "" {
		logger.Info(fmt.Sprintf("✓ Telefone Residencial: %s", telefoneResidencial))
		return telefoneResidencial
	}
	
	// Prioridade 3: Telefone Comercial
	logger.Info("⚠️ Telefone Residencial vazio, tentando Comercial...")
	telefoneComercial, err := ExtractFieldFromTable(ctx, iframeNode, "Telefone Comercial:")
	if err == nil && telefoneComercial != "" {
		logger.Info(fmt.Sprintf("✓ Telefone Comercial: %s", telefoneComercial))
		return telefoneComercial
	}
	
	logger.Info("⚠️ Nenhum telefone encontrado")
	return ""
}