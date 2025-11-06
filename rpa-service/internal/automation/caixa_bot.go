package automation

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

const (
	portalURL = "https://habitacao.caixa.gov.br/siopiweb-web/"
	timeout   = 120 * time.Second
)

// CaixaBot - estrutura principal do bot
type CaixaBot struct{}

// NewCaixaBot - cria uma nova instância do bot
func NewCaixaBot() *CaixaBot {
	return &CaixaBot{}
}

// createBrowserContext - cria contexto do navegador (reutilizável)
func (bot *CaixaBot) createBrowserContext() (context.Context, context.CancelFunc) {
	isHeadless := os.Getenv("HEADLESS") != "false"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", isHeadless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel2 := chromedp.NewContext(allocCtx)
	ctx, cancel3 := context.WithTimeout(ctx, timeout)

	cancelAll := func() {
		cancel3()
		cancel2()
		cancel1()
	}

	return ctx, cancelAll
}

// doLogin - executa o login (função auxiliar privada)
func (bot *CaixaBot) doLogin(ctx context.Context, username, password string) error {
	logger.Info("🔐 Executando login...")

	return chromedp.Run(ctx,
		chromedp.Navigate(portalURL),
		chromedp.Sleep(3*time.Second),

		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.SendKeys(`#username`, username, chromedp.ByID),

		chromedp.WaitVisible(`#password`, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),

		chromedp.WaitVisible(`.btn_login`, chromedp.ByQuery),
		chromedp.Click(`.btn_login`, chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("⏳ Aguardando login processar...")
			return nil
		}),

		chromedp.Sleep(5*time.Second), // Reduzido de 8s para 5s

		chromedp.ActionFunc(func(ctx context.Context) error {
			var currentURL string
			chromedp.Location(&currentURL).Do(ctx)
			logger.Info(fmt.Sprintf("✅ Login concluído - URL: %s", currentURL))
			return nil
		}),
	)
}

// waitForIframe - espera o iframe aparecer e retorna (OTIMIZADO - função auxiliar privada)
func (bot *CaixaBot) waitForIframe(ctx context.Context, stepName string) (*cdp.Node, error) {
	logger.Info(fmt.Sprintf("🎯 [%s] Aguardando iframe...", stepName))

	var iframeNode *cdp.Node

	err := chromedp.Run(ctx,
		// Espera só 2 segundos (otimizado)
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

// fillAndSearchCPF - preenche CPF e clica em buscar (função auxiliar privada)
func (bot *CaixaBot) fillAndSearchCPF(ctx context.Context, cpf string) error {
	logger.Info(fmt.Sprintf("🔍 Preenchendo e buscando CPF: %s", cpf))

	// BUSCA O IFRAME DESTA PÁGINA
	iframeNode, err := bot.waitForIframe(ctx, "Busca CPF")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🎯 Procurando campo #cpfCnpj...")
			return nil
		}),

		// Espera o campo aparecer dentro do iframe
		chromedp.WaitVisible(`#cpfCnpj`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Campo CPF encontrado!")
			return nil
		}),

		chromedp.Click(`#cpfCnpj`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),
		chromedp.Sleep(200*time.Millisecond),

		chromedp.Clear(`#cpfCnpj`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),
		chromedp.Sleep(200*time.Millisecond),

		chromedp.SendKeys(`#cpfCnpj`, cpf, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("✓ CPF digitado: %s", cpf))
			return nil
		}),

		chromedp.Sleep(300*time.Millisecond),

		chromedp.Click(`a[onclick*="executaConsulta"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão Buscar clicado!")
			return nil
		}),

		chromedp.Sleep(4*time.Second), // Reduzido de 6s para 4s
	)
}

// clickProposta - clica na proposta encontrada (função auxiliar privada)
func (bot *CaixaBot) clickProposta(ctx context.Context) error {
	logger.Info("🎯 Procurando proposta para clicar...")

	// BUSCA O IFRAME DESTA NOVA PÁGINA
	iframeNode, err := bot.waitForIframe(ctx, "Lista Propostas")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		// Procura o link da proposta dentro do iframe DESTA página
		chromedp.WaitVisible(`a[onclick*="localizarProposta.do"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Proposta encontrada!")
			return nil
		}),

		chromedp.Click(`a[onclick*="localizarProposta.do"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Proposta clicada! Aguardando detalhes...")
			return nil
		}),

		chromedp.Sleep(4*time.Second), // Reduzido de 5s para 4s
	)
}

// extractAgendamento - extrai a data de agendamento da assinatura (função auxiliar privada)
func (bot *CaixaBot) extractAgendamento(ctx context.Context) (string, error) {
	logger.Info("📊 Extraindo data de agendamento...")

	// BUSCA O IFRAME DESTA NOVA PÁGINA DE DETALHES
	iframeNode, err := bot.waitForIframe(ctx, "Detalhes Proposta")
	if err != nil {
		return "", err
	}

	var agendamento string

	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando 'Agendamento da Assinatura'...")
			
			// XPath dentro do iframe DESTA página
			xpath := `//td[contains(., 'Agendamento da Assinatura')]/following-sibling::td[@class='alinha_esquerda']`
			
			var agendamentoNode []*cdp.Node
			err := chromedp.Nodes(xpath, &agendamentoNode, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao buscar agendamento: %v", err))
				return err
			}

			if len(agendamentoNode) == 0 {
				logger.Error("❌ Data de agendamento não encontrada!")
				return fmt.Errorf("data de agendamento não encontrada")
			}

			// Extrai o texto
			err = chromedp.Text(xpath, &agendamento, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				return err
			}

			logger.Info(fmt.Sprintf("✓ Agendamento extraído: %s", agendamento))
			return nil
		}),
	)

	return agendamento, err
}

// LoginAndSearch - faz login e busca por CPF em uma única operação (FUNÇÃO PRINCIPAL)
func (bot *CaixaBot) LoginAndSearch(req models.LoginAndSearchRequest) (*models.SearchResponse, error) {
	logger.Info("========================================")
	logger.Info("🚀 Iniciando processo: Login + Busca")
	logger.Info(fmt.Sprintf("👤 Usuário: %s", req.Username))
	logger.Info(fmt.Sprintf("📋 CPF: %s", req.CPF))
	logger.Info("========================================")

	ctx, cancel := bot.createBrowserContext()
	defer cancel()

	// 1. Login
	if err := bot.doLogin(ctx, req.Username, req.Password); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro no login: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro no login: %v", err),
		}, err
	}

	// 2. Busca CPF (busca o iframe internamente)
	if err := bot.fillAndSearchCPF(ctx, req.CPF); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao buscar CPF: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao buscar CPF: %v", err),
		}, err
	}

	// 3. Clica na proposta (busca o NOVO iframe internamente)
	if err := bot.clickProposta(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao clicar na proposta: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar na proposta: %v", err),
		}, err
	}

	// 4. Extrai dados (busca o NOVO iframe internamente)
	var clientData models.ClientData
	clientData.CPF = req.CPF

	agendamento, err := bot.extractAgendamento(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair agendamento: %v", err))
		clientData.AgendamentoAssinatura = "Não encontrado"
	} else {
		clientData.AgendamentoAssinatura = agendamento
	}

	clientData.Nome = "Nome do Cliente (A EXTRAIR)"
	clientData.Endereco = "Endereço (A EXTRAIR)"

	logger.Info("========================================")
	logger.Info("✅ PROCESSO CONCLUÍDO!")
	logger.Info(fmt.Sprintf("📝 Nome: %s", clientData.Nome))
	logger.Info(fmt.Sprintf("🏠 Endereço: %s", clientData.Endereco))
	logger.Info(fmt.Sprintf("📅 Agendamento: %s", clientData.AgendamentoAssinatura))
	logger.Info("========================================")

	return &models.SearchResponse{
		Success: true,
		Message: "Dados extraídos com sucesso",
		Data:    &clientData,
	}, nil
}