package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// waitForIframe - espera o iframe aparecer e retorna
func (bot *CaixaBot) waitForIframe(ctx context.Context, stepName string) (*cdp.Node, error) {
	logger.Info(fmt.Sprintf("🎯 [%s] Aguardando iframe...", stepName))

	var iframeNode *cdp.Node

	err := chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// Tenta buscar o iframe até 5 vezes (max 10 segundos)
			for i := 0; i < 5; i++ {
				var iframeNodes []*cdp.Node
				err := chromedp.Nodes(`iframe`, &iframeNodes, chromedp.ByQuery).Do(ctx)
				
				if err == nil && len(iframeNodes) > 0 {
					iframeNode = iframeNodes[0]
					logger.Info(fmt.Sprintf("✓ [%s] Iframe encontrado! (tentativa %d)", stepName, i+1))
					return nil
				}
				
				logger.Info(fmt.Sprintf("⏳ [%s] Iframe não encontrado, tentando novamente... (%d/5)", stepName, i+1))
				time.Sleep(2 * time.Second)
			}
			
			return fmt.Errorf("iframe não encontrado após 5 tentativas")
		}),
	)

	if err != nil {
		logger.Error(fmt.Sprintf("❌ [%s] Erro ao buscar iframe: %v", stepName, err))
		return nil, err
	}

	return iframeNode, nil
}