package navigation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/config"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// MenuNavigator - interface para navegação de menu
// MenuNavigator - interface para navegação de menu
type MenuNavigator interface {
	ClickIrPara(ctx context.Context, iframeWaiter IframeWaiter) error
	ClickMenuOption(ctx context.Context, iframeWaiter IframeWaiter, menuName, optionID string) error
	ClickMenuOptionDirect(ctx context.Context, iframeWaiter IframeWaiter, menuName, optionID string) error  // NOVA
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

// ClickMenuOption - clica em uma opção do menu "Ir para"
func (nav *CaixaMenuNavigator) ClickMenuOption(ctx context.Context, iframeWaiter IframeWaiter, menuName, optionID string) error {
	logger.Info(fmt.Sprintf("🏠 Clicando no menu '%s'...", menuName))
	
	// Busca iframe
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, fmt.Sprintf("Menu %s", menuName))
	if err != nil {
		logger.Error("❌ Iframe não encontrado!")
		return err
	}
	
	logger.Info("✅ Iframe encontrado! Procurando opção do menu...")
	
	// Lista de seletores possíveis baseado no optionID (ex: "imovelPI")
	selectors := []string{
		fmt.Sprintf("#%sDesabCheck", optionID),  // #imovelPIDesabCheck
		fmt.Sprintf("#%s", optionID),            // #imovelPI
		fmt.Sprintf("#%sCheck", optionID),       // #imovelPICheck
		fmt.Sprintf("#%sDesab", optionID),       // #imovelPIDesab
	}
	
	// Tenta cada seletor
	for _, selector := range selectors {
		logger.Info(fmt.Sprintf("🔍 Tentando seletor: %s", selector))
		
		err := chromedp.Run(ctx,
			chromedp.Sleep(1*time.Second),
			chromedp.Click(selector, chromedp.ByID, chromedp.FromNode(iframeNode)),
		)
		
		if err == nil {
			logger.Info(fmt.Sprintf("✅ Menu '%s' clicado: %s", menuName, selector))
			time.Sleep(nav.timeouts.AfterClick)
			return nil
		}
		
		logger.Info(fmt.Sprintf("⚠️ Seletor %s não funcionou", selector))
	}
	
	logger.Error(fmt.Sprintf("❌ Menu '%s' não encontrado!", menuName))
	return fmt.Errorf("menu '%s' não encontrado", menuName)
}

// ClickMenuOptionDirect - clica direto em uma opção do menu (sem abrir "Ir para" antes)
func (nav *CaixaMenuNavigator) ClickMenuOptionDirect(ctx context.Context, iframeWaiter IframeWaiter, menuName, optionID string) error {
	logger.Info(fmt.Sprintf("🏠 Clicando direto no menu '%s'...", menuName))
	
	// Busca iframe
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, fmt.Sprintf("Menu %s", menuName))
	if err != nil {
		logger.Error("❌ Iframe não encontrado!")
		return err
	}
	
	logger.Info("✅ Iframe encontrado! Procurando opção do menu...")
	
	// Lista de seletores possíveis baseado no optionID (ex: "valOperacaoPI")
	selectors := []string{
		fmt.Sprintf("#%sDesabCheck", optionID),  // #valOperacaoPIDesabCheck
		fmt.Sprintf("#%s", optionID),            // #valOperacaoPI
		fmt.Sprintf("#%sCheck", optionID),       // #valOperacaoPICheck
		fmt.Sprintf("#%sDesab", optionID),       // #valOperacaoPIDesab
	}
	
	// Tenta cada seletor
	for _, selector := range selectors {
		logger.Info(fmt.Sprintf("🔍 Tentando seletor: %s", selector))
		
		err := chromedp.Run(ctx,
			chromedp.Sleep(1*time.Second),
			chromedp.Click(selector, chromedp.ByID, chromedp.FromNode(iframeNode)),
		)
		
		if err == nil {
			logger.Info(fmt.Sprintf("✅ Menu '%s' clicado: %s", menuName, selector))
			time.Sleep(nav.timeouts.AfterClick)
			return nil
		}
		
		logger.Info(fmt.Sprintf("⚠️ Seletor %s não funcionou", selector))
	}
	
	logger.Error(fmt.Sprintf("❌ Menu '%s' não encontrado!", menuName))
	return fmt.Errorf("menu '%s' não encontrado", menuName)
}