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

// clickProposta - clica na proposta encontrada
func (bot *CaixaBot) clickProposta(ctx context.Context) error {
	logger.Info("🎯 Procurando proposta para clicar...")

	// Busca o iframe desta nova página
	iframeNode, err := bot.waitForIframe(ctx, "Lista Propostas")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		// Procura o link da proposta dentro do iframe
		chromedp.WaitVisible(`a[onclick*="localizarProposta.do"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Proposta encontrada!")
			return nil
		}),

		chromedp.Click(`a[onclick*="localizarProposta.do"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Proposta clicada! Aguardando detalhes...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}

// extractAgendamento - extrai a data de agendamento da assinatura
func (bot *CaixaBot) extractAgendamento(ctx context.Context) (string, error) {
	logger.Info("📊 Extraindo data de agendamento...")

	// Busca o iframe desta nova página de detalhes
	iframeNode, err := bot.waitForIframe(ctx, "Detalhes Proposta")
	if err != nil {
		return "", err
	}

	var agendamento string

	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando 'Agendamento da Assinatura'...")
			
			xpath := `//td[contains(., 'Agendamento da Assinatura')]/following-sibling::td[@class='alinha_esquerda']`
			
			var agendamentoNode []*cdp.Node
			err := chromedp.Nodes(xpath, &agendamentoNode, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao buscar agendamento: %v", err))
				return err
			}

			if len(agendamentoNode) == 0 {
				logger.Error("❌ Data de agendamento não encontrada!")
				return fmt.Errorf("data de agendamento não encontrada")
			}

			// Extrai o texto
			err = chromedp.Text(xpath, &agendamento, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				return err
			}

			logger.Info(fmt.Sprintf("✓ Agendamento extraído: %s", agendamento))
			return nil
		}),
	)

	return strings.TrimSpace(agendamento), err
}