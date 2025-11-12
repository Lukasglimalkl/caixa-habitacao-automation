package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
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

		chromedp.Sleep(3*time.Second),
	)
}

// waitForMenuDialog - espera o dialog do menu aparecer
func (bot *CaixaBot) waitForMenuDialog(ctx context.Context) error {
	logger.Info("⏳ Aguardando dialog do menu aparecer...")
	
	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Aguarda o dialog aparecer
			for i := 0; i < 5; i++ {
				var dialogNodes []*cdp.Node
				err := chromedp.Nodes(`#divFluxogramaProposta`, &dialogNodes, chromedp.ByID).Do(ctx)
				
				if err == nil && len(dialogNodes) > 0 {
					logger.Info(fmt.Sprintf("✓ Dialog encontrado! (tentativa %d)", i+1))
					return nil
				}
				
				logger.Info(fmt.Sprintf("⏳ Dialog não encontrado, aguardando... (%d/5)", i+1))
				time.Sleep(1 * time.Second)
			}
			
			return fmt.Errorf("dialog não encontrado após 5 tentativas")
		}),
	)
}

// clickMenuImovel - clica no menu "Imóvel" (dentro do dialog que pode ter iframe)
func (bot *CaixaBot) clickMenuImovel(ctx context.Context) error {
	logger.Info("🏠 Clicando no menu 'Imóvel'...")

	return chromedp.Run(ctx,
		// Aguarda o dialog aparecer
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.waitForMenuDialog(ctx)
		}),
		
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando menu 'Imóvel'...")
			
			// ESTRATÉGIA 1: Tenta direto no dialog (sem iframe)
			logger.Info("1️⃣ Tentando encontrar no dialog direto...")
			var nodesDirect []*cdp.Node
			err := chromedp.Nodes(`#imovelPIDesabCheck`, &nodesDirect, chromedp.ByID).Do(ctx)
			
			if err == nil && len(nodesDirect) > 0 {
				logger.Info("✓ Menu encontrado direto no dialog!")
				return chromedp.Click(`#imovelPIDesabCheck`, chromedp.ByID).Do(ctx)
			}
			
			// ESTRATÉGIA 2: Procura iframe DENTRO do dialog
			logger.Info("2️⃣ Procurando iframe dentro do dialog...")
			var dialogIframes []*cdp.Node
			err = chromedp.Nodes(`#divFluxogramaProposta iframe`, &dialogIframes, chromedp.ByQueryAll).Do(ctx)
			
			if err != nil || len(dialogIframes) == 0 {
				logger.Info("⚠️ Nenhum iframe encontrado no dialog")
				
				// ESTRATÉGIA 3: Usa XPath mais agressivo
				logger.Info("3️⃣ Usando XPath para encontrar div com 'imovel'...")
				
				xpathImovel := `//div[contains(@id, 'imovel') and contains(@onclick, 'chamarImovel')]`
				
				var imovelNode []*cdp.Node
				err = chromedp.Nodes(xpathImovel, &imovelNode, chromedp.BySearch).Do(ctx)
				
				if err != nil || len(imovelNode) == 0 {
					logger.Error("❌ Menu Imóvel não encontrado em nenhuma estratégia!")
					return fmt.Errorf("menu imóvel não encontrado")
				}
				
				logger.Info("✓ Menu encontrado via XPath!")
				return chromedp.Click(xpathImovel, chromedp.BySearch).Do(ctx)
			}
			
			// Achou iframe dentro do dialog
			logger.Info(fmt.Sprintf("✓ Iframe encontrado no dialog! Total: %d", len(dialogIframes)))
			
			dialogIframeNode := dialogIframes[0]
			
			// Procura o menu dentro do iframe do dialog
			logger.Info("🔍 Procurando menu dentro do iframe do dialog...")
			
			err = chromedp.WaitVisible(`#imovelPIDesabCheck`, chromedp.ByID, chromedp.FromNode(dialogIframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("❌ Menu não visível no iframe: %v", err))
				return err
			}
			
			logger.Info("✓ Menu 'Imóvel' encontrado no iframe do dialog!")
			
			// Clica no menu dentro do iframe
			return chromedp.Click(`#imovelPIDesabCheck`, chromedp.ByID, chromedp.FromNode(dialogIframeNode)).Do(ctx)
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Menu 'Imóvel' clicado! Aguardando página carregar...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}