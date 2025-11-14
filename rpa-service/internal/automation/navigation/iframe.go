package navigation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/config"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// IframeWaiter - interface para esperar por iframes
type IframeWaiter interface {
	WaitForIframe(ctx context.Context, pageName string) (*cdp.Node, error)
}

// DefaultIframeWaiter - implementação padrão
type DefaultIframeWaiter struct {
	maxRetries int
	waitTime   time.Duration
}

// NewIframeWaiter - cria novo waiter com configurações
func NewIframeWaiter(maxRetries config.MaxRetries, timeouts config.Timeouts) *DefaultIframeWaiter {
	return &DefaultIframeWaiter{
		maxRetries: maxRetries.IframeSearch,
		waitTime:   timeouts.IframeWait / time.Duration(maxRetries.IframeSearch),
	}
}

// WaitForIframe - aguarda o iframe aparecer e retorna o node
func (w *DefaultIframeWaiter) WaitForIframe(ctx context.Context, pageName string) (*cdp.Node, error) {
	logger.Info(fmt.Sprintf("⏳ [%s] Procurando iframe...", pageName))
	
	// Aguarda inicial
	time.Sleep(3 * time.Second)
	
	for tentativa := 1; tentativa <= w.maxRetries; tentativa++ {
		var nodes []*cdp.Node
		
		// SELETOR CORRETO: busca por src="blank.jsp"
		err := chromedp.Run(ctx,
			chromedp.Nodes(`iframe[src="blank.jsp"]`, &nodes, chromedp.BySearch, chromedp.AtLeast(0)),
		)
		
		if err == nil && len(nodes) > 0 {
			logger.Info(fmt.Sprintf("✅ [%s] Iframe encontrado na tentativa %d!", pageName, tentativa))
			time.Sleep(2 * time.Second)
			return nodes[0], nil
		}
		
		if tentativa%3 == 0 {
			logger.Info(fmt.Sprintf("⏳ [%s] Tentativa %d/%d...", pageName, tentativa, w.maxRetries))
		}
		
		time.Sleep(w.waitTime)
	}
	
	logger.Error(fmt.Sprintf("❌ [%s] Iframe não encontrado após %d tentativas!", pageName, w.maxRetries))
	return nil, fmt.Errorf("iframe não encontrado após %d tentativas", w.maxRetries)
}

// WaitForIframeWithSelector - aguarda iframe com seletor customizado
func (w *DefaultIframeWaiter) WaitForIframeWithSelector(ctx context.Context, pageName string, selector string) (*cdp.Node, error) {
	logger.Info(fmt.Sprintf("🎯 [%s] Aguardando iframe com seletor: %s", pageName, selector))
	
	for tentativa := 1; tentativa <= w.maxRetries; tentativa++ {
		var nodes []*cdp.Node
		
		err := chromedp.Run(ctx,
			chromedp.Nodes(selector, &nodes, chromedp.BySearch, chromedp.AtLeast(0)),
		)
		
		if err == nil && len(nodes) > 0 {
			logger.Info(fmt.Sprintf("✓ [%s] Iframe encontrado!", pageName))
			time.Sleep(2 * time.Second)
			return nodes[0], nil
		}
		
		time.Sleep(w.waitTime)
	}
	
	return nil, fmt.Errorf("iframe não encontrado")
}