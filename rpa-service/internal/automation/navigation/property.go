package navigation

import (
	"context"

	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/config"
)

// PropertyNavigator - interface para navegação de imóvel
type PropertyNavigator interface {
	NavigateToProperty(ctx context.Context, menuNav MenuNavigator, iframeWaiter IframeWaiter) error
	NavigateToFinancialValues(ctx context.Context, menuNav MenuNavigator, iframeWaiter IframeWaiter) error
}

// CaixaPropertyNavigator - implementação para navegação de imóvel
type CaixaPropertyNavigator struct {
	timeouts   config.Timeouts
	maxRetries config.MaxRetries
}

// NewCaixaPropertyNavigator - cria novo navegador de imóvel
func NewCaixaPropertyNavigator(timeouts config.Timeouts, maxRetries config.MaxRetries) *CaixaPropertyNavigator {
	return &CaixaPropertyNavigator{
		timeouts:   timeouts,
		maxRetries: maxRetries,
	}
}

// NavigateToProperty - navega até a página de imóvel
func (nav *CaixaPropertyNavigator) NavigateToProperty(ctx context.Context, menuNav MenuNavigator, iframeWaiter IframeWaiter) error {
	// Clica em "Ir para" (precisa abrir menu)
	if err := menuNav.ClickIrPara(ctx, iframeWaiter); err != nil {
		return err
	}
	
	// Clica em "Imóvel"
	return menuNav.ClickMenuOption(ctx, iframeWaiter, "Imóvel", "imovelPI")
}

// NavigateToFinancialValues - navega até valores da operação
func (nav *CaixaPropertyNavigator) NavigateToFinancialValues(ctx context.Context, menuNav MenuNavigator, iframeWaiter IframeWaiter) error {
	// Clica DIRETO em "Valores da Operação" (menu já está aberto)
	return menuNav.ClickMenuOptionDirect(ctx, iframeWaiter, "Valores da Operação", "valOperacaoPI")
}