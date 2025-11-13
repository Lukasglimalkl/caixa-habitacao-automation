package automation

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// ImovelData - dados do imóvel
type ImovelData struct {
	EnderecoImovel string
	CEPImovel      string
}

// extractDadosImovel - extrai dados da página de Imóvel
func (bot *CaixaBot) extractDadosImovel(ctx context.Context) (*ImovelData, error) {
	logger.Info("🏠 Extraindo dados da página de Imóvel...")

	// Aguarda página carregar
	time.Sleep(3 * time.Second)

	// Busca o iframe da página de Imóvel
	iframeNode, err := bot.waitForIframe(ctx, "Página Imóvel")
	if err != nil {
		return nil, err
	}

	var imovelData ImovelData

	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		// Extrai Endereço Completo do Imóvel
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Extraindo endereço do imóvel...")
			
			// XPath para pegar o texto dentro do link <a onclick="exibirDetalheEndereco();">
			xpath := `//tr[.//label[contains(., 'Endereço da Unidade Habitacional:')]]//a[@onclick='exibirDetalheEndereco();']`
			
			var enderecoCompleto string
			err := chromedp.Text(xpath, &enderecoCompleto, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("❌ Erro ao extrair endereço: %v", err))
				return err
			}
			
			enderecoCompleto = strings.TrimSpace(enderecoCompleto)
			logger.Info(fmt.Sprintf("📋 Endereço completo extraído: %s", enderecoCompleto))
			
			// Separa endereço e CEP
			// Regex para encontrar o CEP (formato: CEP XX.XXX-XXX)
			cepRegex := regexp.MustCompile(`CEP\s+(\d{2}\.\d{3}-\d{3})`)
			matches := cepRegex.FindStringSubmatch(enderecoCompleto)
			
			if len(matches) > 1 {
				imovelData.CEPImovel = matches[1] // Captura o CEP
				logger.Info(fmt.Sprintf("✓ CEP Imóvel: %s", imovelData.CEPImovel))
				
				// Pega tudo antes de "CEP"
				indexCEP := strings.Index(enderecoCompleto, "CEP")
				if indexCEP > 0 {
					imovelData.EnderecoImovel = strings.TrimSpace(enderecoCompleto[:indexCEP])
					// Remove vírgula ou espaço final
					imovelData.EnderecoImovel = strings.TrimRight(imovelData.EnderecoImovel, ", ")
					logger.Info(fmt.Sprintf("✓ Endereço Imóvel: %s", imovelData.EnderecoImovel))
				}
			} else {
				// Se não encontrar CEP, usa o endereço completo
				imovelData.EnderecoImovel = enderecoCompleto
				logger.Info("⚠️ CEP não encontrado no endereço")
			}
			
			return nil
		}),
	)

	if err != nil {
		return nil, err
	}

	return &imovelData, nil
}