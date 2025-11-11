package automation

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// clickIrPara - clica no botão "Ir para" (ícone no topo)
func (bot *CaixaBot) clickIrPara(ctx context.Context) error {
	logger.Info("🎯 Clicando no botão 'Ir para'...")

	// Busca o iframe da página atual
	iframeNode, err := bot.waitForIframe(ctx, "Botão Ir Para")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando botão 'Ir para'...")
			return nil
		}),

		chromedp.Sleep(1*time.Second),

		// Procura a imagem com onclick que contém "divFluxogramaProposta"
		chromedp.WaitVisible(`img[onclick*="divFluxogramaProposta"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão 'Ir para' encontrado!")
			return nil
		}),

		// Clica na imagem
		chromedp.Click(`img[onclick*="divFluxogramaProposta"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão 'Ir para' clicado! Aguardando menu aparecer...")
			return nil
		}),

		chromedp.Sleep(2*time.Second),
	)
}

// clickMenuImovel - clica no menu "Imóvel" (SEM iframe - menu é direto na página)
func (bot *CaixaBot) clickMenuImovel(ctx context.Context) error {
	logger.Info("🏠 Clicando no menu 'Imóvel'...")

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando menu 'Imóvel'...")
			return nil
		}),

		chromedp.Sleep(1*time.Second),

		// O menu NÃO está em iframe, busca direto na página
		chromedp.WaitVisible(`#imovelPIDesabCheck`, chromedp.ByID),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Menu 'Imóvel' encontrado!")
			return nil
		}),

		// Clica no div
		chromedp.Click(`#imovelPIDesabCheck`, chromedp.ByID),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Menu 'Imóvel' clicado! Aguardando página carregar...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}