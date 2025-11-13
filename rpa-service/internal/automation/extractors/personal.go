package extractors

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/cdp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
)

// CaixaPersonalExtractor - implementação para dados pessoais
type CaixaPersonalExtractor struct{}

// NewPersonalExtractor - cria novo extrator de dados pessoais
func NewPersonalExtractor() *CaixaPersonalExtractor {
	return &CaixaPersonalExtractor{}
}

// ExtractPersonalData - extrai todos os dados pessoais do participante
func (e *CaixaPersonalExtractor) ExtractPersonalData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	// Número do Contrato
	clientData.NumeroContrato = ExtractFieldWithFallback(ctx, iframeNode, "N° do Contrato:", "Número do Contrato")
	
	// CPF
	clientData.CPF = ExtractFieldWithFallback(ctx, iframeNode, "CPF:", "CPF")
	
	// Nome
	clientData.Nome = ExtractFieldWithFallback(ctx, iframeNode, "Nome:", "Nome")
	
	// Ocupação
	clientData.Ocupacao = ExtractFieldWithFallback(ctx, iframeNode, "Ocupação:", "Ocupação")
	
	// Nacionalidade
	clientData.Nacionalidade = ExtractFieldWithFallback(ctx, iframeNode, "Nacionalidade:", "Nacionalidade")
	
	// Tipo de Identificação e Número (RG/CNH)
	e.extractIdentification(ctx, iframeNode, clientData)
	
	return nil
}

// extractIdentification - extrai tipo de identificação e número (RG/CNH)
func (e *CaixaPersonalExtractor) extractIdentification(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) {
	tipoIdentificacao := ExtractFieldWithFallback(ctx, iframeNode, "Tipo de Identificação:", "Tipo de Identificação")
	clientData.TipoIdentificacao = tipoIdentificacao
	
	if tipoIdentificacao != "" {
		numero := ExtractFieldWithFallback(ctx, iframeNode, "Número:", fmt.Sprintf("Número (%s)", tipoIdentificacao))
		clientData.RG = numero
	}
}