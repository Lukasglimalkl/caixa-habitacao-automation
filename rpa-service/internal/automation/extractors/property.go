package extractors

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// PropertyExtractor - interface para extração de dados do imóvel
type PropertyExtractor interface {
	ExtractPropertyData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error
}

// CaixaPropertyExtractor - implementação para extração de imóvel
type CaixaPropertyExtractor struct{}

// NewPropertyExtractor - cria novo extrator de imóvel
func NewPropertyExtractor() *CaixaPropertyExtractor {
	return &CaixaPropertyExtractor{}
}

// ExtractPropertyData - extrai dados do imóvel
func (e *CaixaPropertyExtractor) ExtractPropertyData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🏠 Extraindo dados do Imóvel...")
	
	time.Sleep(3 * time.Second)
	
	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			return e.extractEnderecoImovel(ctx, iframeNode, clientData)
		}),
	)
}

// extractEnderecoImovel - extrai endereço completo do imóvel
func (e *CaixaPropertyExtractor) extractEnderecoImovel(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🔍 Extraindo endereço do imóvel...")
	
	// XPath para o link que contém o endereço
	xpath := `//tr[.//label[contains(., 'Endereço da Unidade Habitacional:')]]//a[@onclick='exibirDetalheEndereco();']`
	
	var enderecoCompleto string
	err := chromedp.Text(xpath, &enderecoCompleto, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao extrair endereço: %v", err))
		return err
	}
	
	enderecoCompleto = strings.TrimSpace(enderecoCompleto)
	logger.Info(fmt.Sprintf("📋 Endereço completo: %s", enderecoCompleto))
	
	// Separa endereço e CEP
	e.parseEnderecoCompleto(enderecoCompleto, clientData)
	
	return nil
}

// parseEnderecoCompleto - separa endereço e CEP
func (e *CaixaPropertyExtractor) parseEnderecoCompleto(enderecoCompleto string, clientData *models.ClientData) {
	// Regex para encontrar CEP (formato: CEP XX.XXX-XXX)
	cepRegex := regexp.MustCompile(`CEP\s+(\d{2}\.\d{3}-\d{3})`)
	matches := cepRegex.FindStringSubmatch(enderecoCompleto)
	
	if len(matches) > 1 {
		clientData.CEPImovel = matches[1]
		logger.Info(fmt.Sprintf("✓ CEP Imóvel: %s", clientData.CEPImovel))
		
		// Pega tudo antes de "CEP"
		indexCEP := strings.Index(enderecoCompleto, "CEP")
		if indexCEP > 0 {
			clientData.EnderecoImovel = strings.TrimSpace(enderecoCompleto[:indexCEP])
			clientData.EnderecoImovel = strings.TrimRight(clientData.EnderecoImovel, ", ")
			logger.Info(fmt.Sprintf("✓ Endereço Imóvel: %s", clientData.EnderecoImovel))
		}
	} else {
		// Se não encontrar CEP, usa o endereço completo
		clientData.EnderecoImovel = enderecoCompleto
		logger.Info("⚠️ CEP não encontrado no endereço")
	}
}