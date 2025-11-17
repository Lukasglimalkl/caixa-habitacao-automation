package automation

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/config"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
)

// CaixaBot - Robô de automação da Caixa
type CaixaBot struct {
	browserConfig config.BrowserConfig
	timeouts      config.Timeouts
	maxRetries    config.MaxRetries
}

// NewCaixaBot - cria uma nova instância do bot
func NewCaixaBot(headless bool) *CaixaBot {
	return &CaixaBot{
		browserConfig: config.DefaultBrowserConfig(headless),
		timeouts:      config.DefaultTimeouts(),
		maxRetries:    config.DefaultMaxRetries(),
	}
}

// NewCaixaBotWithConfig - cria bot com configuração customizada
func NewCaixaBotWithConfig(browserConfig config.BrowserConfig, timeouts config.Timeouts, maxRetries config.MaxRetries) *CaixaBot {
	return &CaixaBot{
		browserConfig: browserConfig,
		timeouts:      timeouts,
		maxRetries:    maxRetries,
	}
}

//LoginAndSearch - executa login e busca (método principal)
func (bot *CaixaBot) LoginAndSearch(username, password, cpf string) (*models.SearchResponse, error) {
	// Cria contexto base
	ctx := context.Background()
	
	// IMPORTANTE: Cria contexto do Chrome
	browserCtx, cancel := bot.createBrowserContext(ctx)
	defer cancel()
	
	// Cria orquestrador
	orchestrator := NewOrchestrator(bot)
	
	// Executa fluxo completo com o contexto do Chrome
	clientData, err := orchestrator.Execute(browserCtx, username, password, cpf)
	
	if err != nil {
		return &models.SearchResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}
	
	return &models.SearchResponse{
		Success: true,
		Message: "Dados extraídos com sucesso",
		Data:    clientData,
	}, nil
}
// createBrowserContext - cria contexto do navegador
func (bot *CaixaBot) createBrowserContext(ctx context.Context) (context.Context, context.CancelFunc) {
	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, bot.browserConfig.Options...)
	browserCtx, browserCancel := chromedp.NewContext(allocCtx)
	
	// Wrapper para cancelar ambos
	cancelFunc := func() {
		browserCancel()
		allocCancel()
	}
	
	return browserCtx, cancelFunc
}

// GetTimeouts - retorna configurações de timeout
func (bot *CaixaBot) GetTimeouts() config.Timeouts {
	return bot.timeouts
}

// GetMaxRetries - retorna configurações de tentativas
func (bot *CaixaBot) GetMaxRetries() config.MaxRetries {
	return bot.maxRetries
}