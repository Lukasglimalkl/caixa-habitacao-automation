package extractors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// FinancialExtractor - interface para extração de dados financeiros
type FinancialExtractor interface {
	ExtractFinancialData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error
}

// CaixaFinancialExtractor - implementação para extração de valores
type CaixaFinancialExtractor struct{}

// NewFinancialExtractor - cria novo extrator financeiro
func NewFinancialExtractor() *CaixaFinancialExtractor {
	return &CaixaFinancialExtractor{}
}

// ExtractFinancialData - extrai dados financeiros (valor compra e venda)
func (e *CaixaFinancialExtractor) ExtractFinancialData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("💰 Extraindo Valor de Compra e Venda...")
	
	time.Sleep(2 * time.Second)
	
	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			return e.extractValorCompraVenda(ctx, iframeNode, clientData)
		}),
	)
}

// extractValorCompraVenda - extrai valor de compra e venda
func (e *CaixaFinancialExtractor) extractValorCompraVenda(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🔍 Procurando 'Valor Compra e Venda'...")
	
	// XPath para o valor
	xpath := `//tr[.//label[contains(., 'Valor Compra e Venda ou Orçamento Proposto pelo Cliente:')]]//td[@class='alinha_esquerda']`
	
	var valor string
	err := chromedp.Text(xpath, &valor, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao extrair valor: %v", err))
		return err
	}
	
	clientData.ValorCompraVenda = strings.TrimSpace(valor)
	logger.Info(fmt.Sprintf("✓ Valor Compra e Venda: %s", clientData.ValorCompraVenda))
	
	return nil
}