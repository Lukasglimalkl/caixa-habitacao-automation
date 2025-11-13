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

// MenuNavigator - interface para navegação de menu
type MenuNavigator interface {
	ClickIrPara(ctx context.Context, iframeWaiter IframeWaiter) error
	ClickMenuOption(ctx context.Context, iframeWaiter IframeWaiter, optionID string, optionName string) error
}

// CaixaMenuNavigator - implementação para menu da Caixa
type CaixaMenuNavigator struct {
	timeouts   config.Timeouts
	maxRetries config.MaxRetries
}

// NewCaixaMenuNavigator - cria novo navegador de menu
func NewCaixaMenuNavigator(timeouts config.Timeouts, maxRetries config.MaxRetries) *CaixaMenuNavigator {
	return &CaixaMenuNavigator{
		timeouts:   timeouts,
		maxRetries: maxRetries,
	}
}

// ClickIrPara - clica no botão "Ir para"
func (nav *CaixaMenuNavigator) ClickIrPara(ctx context.Context, iframeWaiter IframeWaiter) error {
	logger.Info("🎯 Clicando no botão 'Ir para'...")
	
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Botão Ir Para")
	if err != nil {
		return err
	}
	
	// Procura botão "Ir para"
	xpath := `//img[@onclick="jQuery('#divFluxogramaProposta').dialog('open');outrasOP();"]`
	
	return chromedp.Run(ctx,
		chromedp.WaitVisible(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Click(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Sleep(nav.timeouts.AfterClick),
	)
}

// ClickMenuOption - clica em uma opção específica do menu
func (nav *CaixaMenuNavigator) ClickMenuOption(ctx context.Context, iframeWaiter IframeWaiter, optionID string, optionName string) error {
	logger.Info(fmt.Sprintf("🏠 Clicando no menu '%s'...", optionName))
	
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, fmt.Sprintf("Menu %s", optionName))
	if err != nil {
		return err
	}
	
	// Lista de possíveis IDs (com variações Check/Desab)
	possibleIDs := []string{
		optionID,
		optionID + "Check",
		optionID + "Desab",
		optionID + "DesabCheck",
	}
	
	// Tenta cada ID
	for _, id := range possibleIDs {
		var nodes []*cdp.Node
		err := chromedp.Nodes("#"+id, &nodes, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
		
		if err == nil && len(nodes) > 0 {
			logger.Info(fmt.Sprintf("✓ Botão '%s' encontrado: #%s", optionName, id))
			
			err = chromedp.Click("#"+id, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			if err == nil {
				logger.Info(fmt.Sprintf("✓ Botão '%s' clicado!", optionName))
				time.Sleep(nav.timeouts.AfterClick)
				return nil
			}
		}
	}
	
	return fmt.Errorf("botão '%s' não encontrado", optionName)
}