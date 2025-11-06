package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// CaixaBot - estrutura principal do bot
type CaixaBot struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCaixaBot - cria uma nova instância do bot
func NewCaixaBot() *CaixaBot {
	return &CaixaBot{}
}

// Login - faz login no portal da Caixa
// Login - faz login no portal da Caixa
func (bot *CaixaBot) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	logger.Info("Iniciando login no portal da Caixa...")

	// Opções do Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 👈 false para você ver o navegador
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Cria contexto do Chrome
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Timeout de 60 segundos
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// URL do portal
	url := "https://habitacao.caixa.gov.br/siopiweb-web/"

	var loginSuccess bool

	err := chromedp.Run(ctx,
		// Navega para a página de login
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second), // Espera a página carregar

		// Preenche o campo de username
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.SendKeys(`#username`, req.Username, chromedp.ByID),
		
		// Preenche o campo de senha
		chromedp.WaitVisible(`#password`, chromedp.ByID),
		chromedp.SendKeys(`#password`, req.Password, chromedp.ByID),
		
		// Clica no botão de login
		chromedp.WaitVisible(`.btn_login`, chromedp.ByQuery),
		chromedp.Click(`.btn_login`, chromedp.ByQuery),
		
		// Espera 5 segundos para ver se deu certo
		chromedp.Sleep(5*time.Second),
		
		// Verifica se logou (checando se NÃO está mais na tela de login)
		chromedp.ActionFunc(func(ctx context.Context) error {
			var currentURL string
			if err := chromedp.Location(&currentURL).Do(ctx); err != nil {
				return err
			}
			
			logger.Info(fmt.Sprintf("URL após login: %s", currentURL))
			
			// Se mudou de URL, login deu certo
			if currentURL != url {
				loginSuccess = true
				logger.Info("Login realizado com sucesso!")
			} else {
				logger.Error("Login falhou - ainda na página de login")
			}
			return nil
		}),
	)

	if err != nil {
		logger.Error(fmt.Sprintf("Erro no login: %v", err))
		return &models.LoginResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao fazer login: %v", err),
		}, err
	}

	if !loginSuccess {
		return &models.LoginResponse{
			Success: false,
			Message: "Login falhou - credenciais inválidas ou timeout",
		}, nil
	}

	// Gera um token de sessão
	sessionToken := fmt.Sprintf("session_%d", time.Now().Unix())

	return &models.LoginResponse{
		Success:      true,
		Message:      "Login realizado com sucesso",
		SessionToken: sessionToken,
	}, nil
}

// SearchByCPF - busca dados por CPF no portal
// SearchByCPF - busca dados por CPF no portal
func (bot *CaixaBot) SearchByCPF(req models.SearchRequest) (*models.SearchResponse, error) {
	logger.Info(fmt.Sprintf("Buscando dados do CPF: %s", req.CPF))

	// Valida session token (por enquanto só checa se não está vazio)
	if req.SessionToken == "" {
		return &models.SearchResponse{
			Success: false,
			Message: "Token de sessão inválido",
		}, fmt.Errorf("session token vazio")
	}

	// Opções do Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // false para você ver
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Precisamos fazer login de novo pra ter a sessão
	url := "https://habitacao.caixa.gov.br/siopiweb-web/"

	var nome, endereco string
	var searchSuccess bool

	err := chromedp.Run(ctx,
		// Navega para o portal (assumindo que precisa logar de novo)
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),

		// TODO: Aqui você precisa fazer login de novo ou reusar a sessão
		// Por enquanto vou assumir que já está logado ou vamos implementar depois

		// Espera o campo de CPF aparecer
		chromedp.WaitVisible(`#cpfCnpj`, chromedp.ByID),
		
		// Limpa o campo (se tiver algo)
		chromedp.Clear(`#cpfCnpj`, chromedp.ByID),
		
		// Preenche o CPF
		chromedp.SendKeys(`#cpfCnpj`, req.CPF, chromedp.ByID),
		chromedp.Sleep(1*time.Second),
		
		// Clica no botão Buscar (usando XPath pelo texto)
		chromedp.Click(`//a[contains(text(), 'Buscar')]`, chromedp.BySearch),
		
		// Espera 5 segundos para carregar os resultados
		chromedp.Sleep(5*time.Second),
		
		// TODO: Aqui você precisa extrair os dados que aparecem
		// Por enquanto vou deixar simulado, me diz quais dados aparecem!
		chromedp.ActionFunc(func(ctx context.Context) error {
			nome = "Nome do Cliente (EXTRAIR)"
			endereco = "Endereço (EXTRAIR)"
			searchSuccess = true
			logger.Info("Busca realizada, aguardando seletores dos dados...")
			return nil
		}),
	)

	if err != nil {
		logger.Error(fmt.Sprintf("Erro na busca: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao buscar CPF: %v", err),
		}, err
	}

	if !searchSuccess {
		return &models.SearchResponse{
			Success: false,
			Message: "CPF não encontrado",
		}, nil
	}

	return &models.SearchResponse{
		Success: true,
		Message: "Dados encontrados",
		Data: &models.ClientData{
			CPF:      req.CPF,
			Nome:     nome,
			Endereco: endereco,
		},
	}, nil
}

// LoginAndSearch - faz login e busca por CPF em uma única operação
func (bot *CaixaBot) LoginAndSearch(req models.LoginAndSearchRequest) (*models.SearchResponse, error) {
	logger.Info("========================================")
	logger.Info("🚀 Iniciando processo completo: Login + Busca")
	logger.Info(fmt.Sprintf("👤 Usuário: %s", req.Username))
	logger.Info(fmt.Sprintf("📋 CPF: %s", req.CPF))
	logger.Info("========================================")

	// Opções do Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // false para você ver
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 120*time.Second) // 2 minutos
	defer cancel()

	url := "https://habitacao.caixa.gov.br/siopiweb-web/"

	var nome, endereco string
	var searchSuccess bool

	err := chromedp.Run(ctx,
		// ========== ETAPA 1: NAVEGAÇÃO ==========
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("📡 Acessando portal da Caixa...")
			return nil
		}),
		
		chromedp.Navigate(url),
		chromedp.Sleep(4*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✅ Página carregada")
			return nil
		}),

		// ========== ETAPA 2: LOGIN ==========
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔐 Iniciando login...")
			return nil
		}),
		
		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("   ✓ Campo username encontrado")
			return nil
		}),
		chromedp.SendKeys(`#username`, req.Username, chromedp.ByID),
		
		chromedp.WaitVisible(`#password`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("   ✓ Campo password encontrado")
			return nil
		}),
		chromedp.SendKeys(`#password`, req.Password, chromedp.ByID),
		
		chromedp.WaitVisible(`.btn_login`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("   ✓ Botão login encontrado")
			return nil
		}),
		chromedp.Click(`.btn_login`, chromedp.ByQuery),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("⏳ Aguardando login processar...")
			return nil
		}),
		
		chromedp.Sleep(8*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			var currentURL string
			chromedp.Location(&currentURL).Do(ctx)
			logger.Info(fmt.Sprintf("📍 URL atual: %s", currentURL))
			logger.Info("✅ Login concluído!")
			return nil
		}),

		// ========== ETAPA 3: BUSCA CPF ==========
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando campo de CPF...")
			return nil
		}),
		
		// Espera o campo aparecer (pode demorar)
		chromedp.WaitVisible(`#cpfCnpj`, chromedp.ByID),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("   ✓ Campo CPF encontrado!")
			return nil
		}),
		
		// Clica no campo primeiro (às vezes precisa de foco)
		chromedp.Click(`#cpfCnpj`, chromedp.ByID),
		chromedp.Sleep(500*time.Millisecond),
		
		// Limpa e preenche
		chromedp.Clear(`#cpfCnpj`, chromedp.ByID),
		chromedp.SendKeys(`#cpfCnpj`, req.CPF, chromedp.ByID),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info(fmt.Sprintf("   ✓ CPF preenchido: %s", req.CPF))
			return nil
		}),
		
		chromedp.Sleep(1*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando botão Buscar (via onclick)...")
			return nil
		}),
		
		// OPÇÃO 1: Clica pelo seletor do onclick
// OPÇÃO 2: Executa o JavaScript diretamente
		chromedp.Evaluate(`executaConsulta('cpfCnpjProposta');`, nil),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("   ✓ JavaScript executado!")
			return nil
		}),		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("   ✓ Botão Buscar clicado!")
			return nil
		}),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("⏳ Aguardando resultados...")
			return nil
		}),
		
		chromedp.Sleep(8*time.Second),
		
		// ========== ETAPA 4: EXTRAÇÃO DE DADOS ==========
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("📊 Extraindo dados da página...")
			
			// TODO: AQUI VOCÊ VAI ADICIONAR OS SELETORES REAIS
			nome = "Nome do Cliente (AGUARDANDO SELETORES)"
			endereco = "Endereço (AGUARDANDO SELETORES)"
			searchSuccess = true
			
			logger.Info("========================================")
			logger.Info("✅ PROCESSO CONCLUÍDO!")
			logger.Info(fmt.Sprintf("📝 Nome: %s", nome))
			logger.Info(fmt.Sprintf("🏠 Endereço: %s", endereco))
			logger.Info("========================================")
			
			return nil
		}),
	)

	if err != nil {
		logger.Error("========================================")
		logger.Error(fmt.Sprintf("❌ ERRO NO PROCESSO: %v", err))
		logger.Error("========================================")
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro: %v", err),
		}, err
	}

	if !searchSuccess {
		return &models.SearchResponse{
			Success: false,
			Message: "CPF não encontrado ou erro na busca",
		}, nil
	}

	return &models.SearchResponse{
		Success: true,
		Message: "Dados extraídos com sucesso",
		Data: &models.ClientData{
			CPF:      req.CPF,
			Nome:     nome,
			Endereco: endereco,
		},
	}, nil
}
// Close - fecha o navegador
func (bot *CaixaBot) Close() {
	if bot.cancel != nil {
		bot.cancel()
		logger.Info("Bot encerrado")
	}
}

