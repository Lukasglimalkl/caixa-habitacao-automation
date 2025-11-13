package automation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// clickValoresOperacao - clica no menu "Valores da Operação"
func (bot *CaixaBot) clickValoresOperacao(ctx context.Context) error {
	logger.Info("💰 Clicando no menu 'Valores da Operação'...")

	// Busca o iframe da página atual
	iframeNode, err := bot.waitForIframe(ctx, "Valores Operação")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando botão 'Valores da Operação'...")
			
			// Lista de IDs possíveis
			possibleIDs := []string{
				"valOperacaoPIDesabCheck",
				"valOperacaoPI",
				"valOperacaoPICheck",
				"valOperacaoPIDesab",
			}
			
			// Tenta cada ID
			for _, id := range possibleIDs {
				var nodes []*cdp.Node
				err := chromedp.Nodes(`#`+id, &nodes, chromedp.ByID, chromedp.FromNode(iframeNode)).Do(ctx)
				
				if err == nil && len(nodes) > 0 {
					logger.Info(fmt.Sprintf("✓ Botão encontrado: #%s", id))
					return chromedp.Click(`#`+id, chromedp.ByID, chromedp.FromNode(iframeNode)).Do(ctx)
				}
			}
			
			return fmt.Errorf("botão 'Valores da Operação' não encontrado")
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão 'Valores da Operação' clicado! Aguardando página carregar...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}

// extractValorCompraVenda - extrai o valor de compra e venda
func (bot *CaixaBot) extractValorCompraVenda(ctx context.Context) (string, error) {
	logger.Info("💰 Extraindo Valor de Compra e Venda...")

	// Aguarda página carregar
	time.Sleep(2 * time.Second)

	// Busca o iframe da página de Valores
	iframeNode, err := bot.waitForIframe(ctx, "Extração Valor Compra")
	if err != nil {
		return "", err
	}

	var valorCompraVenda string

	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando 'Valor Compra e Venda ou Orçamento Proposto pelo Cliente'...")
			
			// XPath específico para o valor de compra e venda
			xpath := `//tr[.//label[contains(., 'Valor Compra e Venda ou Orçamento Proposto pelo Cliente:')]]//td[@class='alinha_esquerda']`
			
			var valor string
			err := chromedp.Text(xpath, &valor, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("❌ Erro ao extrair valor: %v", err))
				return err
			}
			
			valorCompraVenda = strings.TrimSpace(valor)
			logger.Info(fmt.Sprintf("✓ Valor Compra e Venda: %s", valorCompraVenda))
			
			return nil
		}),
	)

	if err != nil {
		return "", err
	}

	return valorCompraVenda, nil
}