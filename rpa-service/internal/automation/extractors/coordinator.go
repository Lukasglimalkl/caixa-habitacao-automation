package extractors

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// DataCoordinator - coordena todas as extrações
type DataCoordinator struct {
	personalExtractor  *CaixaPersonalExtractor
	contactExtractor   *CaixaContactExtractor
	addressExtractor   *CaixaAddressExtractor
	bankingExtractor   *CaixaBankingExtractor
	propertyExtractor  *CaixaPropertyExtractor
	financialExtractor *CaixaFinancialExtractor
}

// NewDataCoordinator - cria novo coordenador de extração
func NewDataCoordinator() *DataCoordinator {
	return &DataCoordinator{
		personalExtractor:  NewPersonalExtractor(),
		contactExtractor:   NewContactExtractor(),
		addressExtractor:   NewAddressExtractor(),
		bankingExtractor:   NewBankingExtractor(),
		propertyExtractor:  NewPropertyExtractor(),
		financialExtractor: NewFinancialExtractor(),
	}
}

// ExtractAllParticipantData - orquestra extração de dados do participante
func (c *DataCoordinator) ExtractAllParticipantData(ctx context.Context, iframeNode *cdp.Node) (*models.ClientData, error) {
	logger.Info("📊 Iniciando extração completa de dados do participante...")
	
	clientData := &models.ClientData{}
	
	err := chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		
		// 1. Dados Pessoais
		chromedp.ActionFunc(func(ctx context.Context) error {
			return c.personalExtractor.ExtractPersonalData(ctx, iframeNode, clientData)
		}),
		
		// 2. Dados de Contato
		chromedp.ActionFunc(func(ctx context.Context) error {
			return c.contactExtractor.ExtractContactData(ctx, iframeNode, clientData)
		}),
		
		// 3. Dados de Endereço
		chromedp.ActionFunc(func(ctx context.Context) error {
			return c.addressExtractor.ExtractAddressData(ctx, iframeNode, clientData)
		}),
		
		chromedp.Sleep(1*time.Second),
		
		// 4. Dados Bancários
		chromedp.ActionFunc(func(ctx context.Context) error {
			return c.bankingExtractor.ExtractBankingData(ctx, iframeNode, clientData)
		}),
	)
	
	if err != nil {
		logger.Error("❌ Erro durante extração de dados")
		return nil, err
	}
	
	logger.Info("✅ Extração completa de dados do participante finalizada!")
	return clientData, nil
}

// ExtractPropertyData - extrai dados do imóvel
func (c *DataCoordinator) ExtractPropertyData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	return c.propertyExtractor.ExtractPropertyData(ctx, iframeNode, clientData)
}

// ExtractFinancialData - extrai dados financeiros
func (c *DataCoordinator) ExtractFinancialData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	return c.financialExtractor.ExtractFinancialData(ctx, iframeNode, clientData)
}