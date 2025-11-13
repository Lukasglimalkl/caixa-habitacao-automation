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

// ExtractAddressData - extrai todos os dados de endereço
func ExtractAddressData(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🏠 Extraindo dados de endereço...")
	
	// Faz scroll até a tabela de endereço
	ScrollToTable(ctx, "Endereço")
	
	// XPath base: pega a primeira tabela com "Endereço" no header
	baseXPath := `//table[.//th[contains(text(), 'Endereço') and not(contains(text(), 'Correspondência'))]]`
	
	// Extrai cada campo do endereço
	extractCEP(ctx, iframeNode, baseXPath, clientData)
	extractTipoLogradouro(ctx, iframeNode, baseXPath, clientData)
	extractLogradouro(ctx, iframeNode, baseXPath, clientData)
	extractNumeroEndereco(ctx, iframeNode, baseXPath, clientData)
	extractBairro(ctx, iframeNode, baseXPath, clientData)
	extractMunicipioUF(ctx, iframeNode, baseXPath, clientData)
	extractComplemento(ctx, iframeNode, baseXPath, clientData)
	
	return nil
}

func extractCEP(ctx context.Context, iframeNode *cdp.Node, baseXPath string, clientData *models.ClientData) {
	xpath := baseXPath + `//tr[.//label[contains(., 'CEP:')]]/td[@class='alinha_esquerda'][1]`
	var cep string
	err := chromedp.Text(xpath, &cep, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err == nil {
		clientData.CEP = strings.TrimSpace(cep)
		logger.Info(fmt.Sprintf("✓ CEP: %s", clientData.CEP))
	}
}

func extractTipoLogradouro(ctx context.Context, iframeNode *cdp.Node, baseXPath string, clientData *models.ClientData) {
	xpath := baseXPath + `//tr[.//label[contains(., 'Tipo de Logradouro:')]]/td[@class='alinha_esquerda'][last()]`
	var tipoLogradouro string
	err := chromedp.Text(xpath, &tipoLogradouro, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err == nil {
		clientData.TipoLogradouro = strings.TrimSpace(tipoLogradouro)
		logger.Info(fmt.Sprintf("✓ Tipo de Logradouro: %s", clientData.TipoLogradouro))
	}
}

func extractLogradouro(ctx context.Context, iframeNode *cdp.Node, baseXPath string, clientData *models.ClientData) {
	xpath := baseXPath + `//tr[.//label[contains(., 'Logradouro:')]]/td[@class='alinha_esquerda'][1]`
	var logradouro string
	err := chromedp.Text(xpath, &logradouro, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err == nil {
		clientData.Logradouro = strings.TrimSpace(logradouro)
		logger.Info(fmt.Sprintf("✓ Logradouro: %s", clientData.Logradouro))
	}
}

func extractNumeroEndereco(ctx context.Context, iframeNode *cdp.Node, baseXPath string, clientData *models.ClientData) {
	xpath := baseXPath + `//tr[.//label[contains(., 'Número:')]]/td[@class='alinha_esquerda'][last()]`
	var numero string
	err := chromedp.Text(xpath, &numero, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err == nil {
		clientData.Numero = strings.TrimSpace(numero)
		logger.Info(fmt.Sprintf("✓ Número: %s", clientData.Numero))
	}
}

func extractBairro(ctx context.Context, iframeNode *cdp.Node, baseXPath string, clientData *models.ClientData) {
	xpath := baseXPath + `//tr[.//label[contains(., 'Bairro:')]]/td[@class='alinha_esquerda'][last()]`
	var bairro string
	err := chromedp.Text(xpath, &bairro, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err == nil {
		clientData.Bairro = strings.TrimSpace(bairro)
		logger.Info(fmt.Sprintf("✓ Bairro: %s", clientData.Bairro))
	}
}

func extractMunicipioUF(ctx context.Context, iframeNode *cdp.Node, baseXPath string, clientData *models.ClientData) {
	xpath := baseXPath + `//tr[.//label[contains(., 'Município - UF:')]]/td[@class='alinha_esquerda']`
	var municipioUF string
	err := chromedp.Text(xpath, &municipioUF, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err == nil {
		municipioUF = strings.TrimSpace(municipioUF)
		partes := strings.Split(municipioUF, "-")
		if len(partes) >= 2 {
			clientData.Municipio = strings.TrimSpace(partes[0])
			clientData.UF = strings.TrimSpace(partes[1])
		} else {
			clientData.Municipio = municipioUF
		}
		logger.Info(fmt.Sprintf("✓ Município: %s", clientData.Municipio))
		logger.Info(fmt.Sprintf("✓ UF: %s", clientData.UF))
	}
}

func extractComplemento(ctx context.Context, iframeNode *cdp.Node, baseXPath string, clientData *models.ClientData) {
	xpath := baseXPath + `//tr[.//label[contains(., 'Complemento:')]]/td[@class='alinha_esquerda'][1]`
	var complemento string
	chromedp.Text(xpath, &complemento, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	clientData.Complemento = strings.TrimSpace(complemento)
	if clientData.Complemento != "" {
		logger.Info(fmt.Sprintf("✓ Complemento: %s", clientData.Complemento))
	}
}