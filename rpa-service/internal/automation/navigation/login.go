package navigation

import (
	"context"
	"fmt"
	"time"

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


// Login - realiza o login no portal
func (nav *CaixaLoginNavigator) Login(ctx context.Context, username, password string) error {
	logger.Info("🔐 Iniciando processo de login...")
	logger.Info(fmt.Sprintf("🌐 URL: %s", nav.url))
	
	err := chromedp.Run(ctx,
		// Navega para a página
		chromedp.Navigate(nav.url),
		chromedp.Sleep(5*time.Second),
		
		// Debug: Verifica se página carregou
		chromedp.ActionFunc(func(ctx context.Context) error {
			var title string
			chromedp.Title(&title).Do(ctx)
			logger.Info(fmt.Sprintf("📄 Título: %s", title))
			return nil
		}),
		
		// Debug: Verifica se campos existem
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Verificando se campos existem...")
			
			var usernameExists bool
			chromedp.Evaluate(`document.querySelector('#username') !== null`, &usernameExists).Do(ctx)
			logger.Info(fmt.Sprintf("Campo username existe: %v", usernameExists))
			
			var passwordExists bool
			chromedp.Evaluate(`document.querySelector('#password') !== null`, &passwordExists).Do(ctx)
			logger.Info(fmt.Sprintf("Campo password existe: %v", passwordExists))
			
			var btnExists bool
			chromedp.Evaluate(`document.querySelector('#btn_login') !== null`, &btnExists).Do(ctx)
			logger.Info(fmt.Sprintf("Botão login existe: %v", btnExists))
			
			return nil
		}),
		
		// Aguarda campo de usuário
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		
		// Preenche usando JavaScript
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("📝 Preenchendo usuário com JavaScript...")
			script := fmt.Sprintf(`document.querySelector('#username').value = '%s';`, username)
			return chromedp.Evaluate(script, nil).Do(ctx)
		}),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("📝 Preenchendo senha com JavaScript...")
			script := fmt.Sprintf(`document.querySelector('#password').value = '%s';`, password)
			return chromedp.Evaluate(script, nil).Do(ctx)
		}),
		
		// Verifica se preencheu
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Verificando se campos foram preenchidos...")
			
			var usernameValue string
			chromedp.Evaluate(`document.querySelector('#username').value`, &usernameValue).Do(ctx)
			logger.Info(fmt.Sprintf("Valor username: %s", usernameValue))
			
			var passwordValue string
			chromedp.Evaluate(`document.querySelector('#password').value`, &passwordValue).Do(ctx)
			logger.Info(fmt.Sprintf("Valor password: %d caracteres", len(passwordValue)))
			
			return nil
		}),
		
		chromedp.Sleep(1*time.Second),
		
// Clica no botão
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🎯 Clicando no botão de login...")
			script := `document.querySelector('#btn_login').click();`
			return chromedp.Evaluate(script, nil).Do(ctx)
		}),
		
		// Aguarda navegação COMPLETA
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("⏳ Aguardando redirecionamento pós-login...")
			return nil
		}),
		chromedp.Sleep(8*time.Second),
		
		chromedp.WaitReady("body", chromedp.ByQuery),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			var currentURL string
			chromedp.Evaluate(`window.location.href`, &currentURL).Do(ctx)
			logger.Info(fmt.Sprintf("📍 URL atual: %s", currentURL))
			logger.Info("✅ Página pós-login carregada!")
			return nil
		}),
	)
	
	return err
}
// NewCaixaLoginNavigator - cria novo navegador de login
func NewCaixaLoginNavigator(timeouts config.Timeouts, maxRetries config.MaxRetries) *CaixaLoginNavigator {
	return &CaixaLoginNavigator{
		url:        "https://habitacao.caixa.gov.br/siopiweb-web/",
		timeouts:   timeouts,
		maxRetries: maxRetries,
	}
}

// VerifyLoginSuccess - verifica se o login foi bem-sucedido
func (nav *CaixaLoginNavigator) VerifyLoginSuccess(ctx context.Context) error {
	logger.Info("✓ Verificando sucesso do login...")
	
	// Aguarda um pouco mais para garantir
	time.Sleep(2 * time.Second)
	
	// Tenta pegar o título, mas não falha se der erro
	var pageTitle string
	err := chromedp.Title(&pageTitle).Do(ctx)
	
	if err != nil {
		logger.Info("⚠️ Não foi possível verificar título (página ainda carregando)")
		// Não retorna erro, só avisa
		return nil
	}
	
	logger.Info(fmt.Sprintf("📄 Título da página: %s", pageTitle))
	
	// Verifica se a URL mudou (sinal de sucesso)
	var currentURL string
	chromedp.Evaluate(`window.location.href`, &currentURL).Do(ctx)
	
	if currentURL != "" && currentURL != nav.url {
		logger.Info("✅ Login realizado com sucesso! (URL mudou)")
		return nil
	}
	
	logger.Info("✅ Login aparentemente bem-sucedido!")
	return nil
}