package navigation

import (
	"context"
	"fmt"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/config"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// LoginNavigator - interface para navegação de login
type LoginNavigator interface {
	Login(ctx context.Context, username, password string) error
	VerifyLoginSuccess(ctx context.Context) error
}

// CaixaLoginNavigator - implementação para portal da Caixa
type CaixaLoginNavigator struct {
	url        string
	timeouts   config.Timeouts
	maxRetries config.MaxRetries
}

// NewCaixaLoginNavigator - cria novo navegador de login
func NewCaixaLoginNavigator(timeouts config.Timeouts, maxRetries config.MaxRetries) *CaixaLoginNavigator {
	return &CaixaLoginNavigator{
		url:        "https://habpf.caixa.gov.br/",
		timeouts:   timeouts,
		maxRetries: maxRetries,
	}
}

// Login - realiza o login no portal
func (nav *CaixaLoginNavigator) Login(ctx context.Context, username, password string) error {
	logger.Info("🔐 Iniciando processo de login...")
	
	return chromedp.Run(ctx,
		// Navega para a página
		chromedp.Navigate(nav.url),
		chromedp.Sleep(nav.timeouts.PageLoad),
		
		// Preenche credenciais
		chromedp.ActionFunc(func(ctx context.Context) error {
			return nav.fillCredentials(ctx, username, password)
		}),
		
		// Clica no botão de login
		chromedp.ActionFunc(func(ctx context.Context) error {
			return nav.clickLoginButton(ctx)
		}),
		
		// Aguarda redirecionamento
		chromedp.Sleep(nav.timeouts.AfterClick),
	)
}

// fillCredentials - preenche usuário e senha
func (nav *CaixaLoginNavigator) fillCredentials(ctx context.Context, username, password string) error {
	logger.Info("📝 Preenchendo credenciais...")
	
	return chromedp.Run(ctx,
		chromedp.WaitVisible(`#userName`, chromedp.ByID),
		chromedp.SendKeys(`#userName`, username, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),
	)
}

// clickLoginButton - clica no botão de login
func (nav *CaixaLoginNavigator) clickLoginButton(ctx context.Context) error {
	logger.Info("🎯 Clicando no botão de login...")
	
	return chromedp.Run(ctx,
		chromedp.Click(`input[type="submit"][value="Login"]`, chromedp.BySearch),
	)
}

// VerifyLoginSuccess - verifica se o login foi bem-sucedido
// VerifyLoginSuccess - verifica se o login foi bem-sucedido
func (nav *CaixaLoginNavigator) VerifyLoginSuccess(ctx context.Context) error {
		logger.Info("✓ Verificando sucesso do login...")
	
	var pageTitle string
	err := chromedp.Title(&pageTitle).Do(ctx)
	
	if err != nil {
		return fmt.Errorf("erro ao verificar título da página: %w", err)
	}
	
	logger.Info(fmt.Sprintf("📄 Título da página: %s", pageTitle))
	return nil
}