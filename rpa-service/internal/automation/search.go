package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// fillAndSearchCPF - preenche CPF e clica em buscar
func (bot *CaixaBot) fillAndSearchCPF(ctx context.Context, cpf string) error {
	logger.Info(fmt.Sprintf("🔍 Preenchendo e buscando CPF: %s", cpf))

	// Busca o iframe desta página
	iframeNode, err := bot.waitForIframe(ctx, "Busca CPF")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🎯 Procurando campo #cpfCnpj...")
			return nil
		}),

		// Espera o campo aparecer dentro do iframe
		chromedp.WaitVisible(`#cpfCnpj`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Campo CPF encontrado!")
			return nil
		}),

		chromedp.Click(`#cpfCnpj`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),
		chromedp.Sleep(200*time.Millisecond),

		chromedp.Clear(`#cpfCnpj`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),
		chromedp.Sleep(200*time.Millisecond),

		chromedp.SendKeys(`#cpfCnpj`, cpf, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("✓ CPF digitado: %s", cpf))
			return nil
		}),

		chromedp.Sleep(300*time.Millisecond),

		chromedp.Click(`a[onclick*="executaConsulta"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão Buscar clicado!")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}