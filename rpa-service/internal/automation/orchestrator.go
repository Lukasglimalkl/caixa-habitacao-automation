package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/extractors"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/navigation"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// Orchestrator - orquestra todo o fluxo de automação
type Orchestrator struct {
	bot              *CaixaBot
	iframeWaiter     navigation.IframeWaiter
	loginNav         navigation.LoginNavigator
	searchNav        navigation.SearchNavigator
	participantsNav  navigation.ParticipantsNavigator
	menuNav          navigation.MenuNavigator
	propertyNav      navigation.PropertyNavigator
	dataCoordinator  *extractors.DataCoordinator
}

// NewOrchestrator - cria novo orquestrador
func NewOrchestrator(bot *CaixaBot) *Orchestrator {
	timeouts := bot.GetTimeouts()
	maxRetries := bot.GetMaxRetries()
	
	return &Orchestrator{
		bot:              bot,
		iframeWaiter:     navigation.NewIframeWaiter(maxRetries, timeouts),
		loginNav:         navigation.NewCaixaLoginNavigator(timeouts, maxRetries),
		searchNav:        navigation.NewCaixaSearchNavigator(timeouts, maxRetries),
		participantsNav:  navigation.NewCaixaParticipantsNavigator(timeouts, maxRetries),
		menuNav:          navigation.NewCaixaMenuNavigator(timeouts, maxRetries),
		propertyNav:      navigation.NewCaixaPropertyNavigator(timeouts, maxRetries),
		dataCoordinator:  extractors.NewDataCoordinator(),
	}
}

// Execute - executa o fluxo completo de automação
func (o *Orchestrator) Execute(ctx context.Context, username, password, cpf string) (*models.ClientData, error) {
	logger.Info("🚀 Iniciando processo de automação completo...")
	logger.Info("========================================")
	
	// Cria context do navegador
	browserCtx, cancel := o.bot.createBrowserContext(ctx)
	defer cancel()
	
	// Context com timeout
	timeoutCtx, timeoutCancel := context.WithTimeout(browserCtx, o.bot.timeouts.BrowserContext)
	defer timeoutCancel()
	
	// Inicia navegador
	if err := chromedp.Run(timeoutCtx); err != nil {
		return nil, fmt.Errorf("erro ao iniciar navegador: %w", err)
	}
	
	// 1. Login
	if err := o.executeLogin(timeoutCtx, username, password); err != nil {
		return nil, fmt.Errorf("erro no login: %w", err)
	}
	
	// 2. Busca por CPF
	if err := o.executeSearch(timeoutCtx, cpf); err != nil {
		return nil, fmt.Errorf("erro na busca: %w", err)
	}
	
	// 3. Extrai dados do participante
	clientData, err := o.extractParticipantData(timeoutCtx)
	if err != nil {
		return nil, fmt.Errorf("erro ao extrair dados: %w", err)
	}
	
	// 4. Navega e extrai dados do imóvel
	if err := o.extractPropertyData(timeoutCtx, clientData); err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair dados do imóvel: %v", err))
		// Não falha o processo
	}
	
	// 5. Navega e extrai dados financeiros
	if err := o.extractFinancialData(timeoutCtx, clientData); err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair dados financeiros: %v", err))
		// Não falha o processo
	}
	
	// Log final
	o.logFinalResults(clientData)
	
	return clientData, nil
}

// executeLogin - executa o processo de login
func (o *Orchestrator) executeLogin(ctx context.Context, username, password string) error {
	logger.Info("========================================")
	logger.Info("ETAPA 1: LOGIN")
	logger.Info("========================================")
	
	if err := o.loginNav.Login(ctx, username, password); err != nil {
		return err
	}
	
	if err := o.loginNav.VerifyLoginSuccess(ctx); err != nil {
		return err
	}
	
	logger.Info("✅ Login realizado com sucesso!")
	return nil
}

// executeSearch - executa a busca por CPF
func (o *Orchestrator) executeSearch(ctx context.Context, cpf string) error {
	logger.Info("========================================")
	logger.Info("ETAPA 2: BUSCA POR CPF")
	logger.Info("========================================")
	
	if err := o.searchNav.SearchByCPF(ctx, cpf); err != nil {
		return err
	}
	
	if err := o.searchNav.ClickFirstResult(ctx); err != nil {
		return err
	}
	
	// Extrai agendamento de assinatura
	agendamento, err := o.searchNav.ExtractAgendamentoAssinatura(ctx, o.iframeWaiter)
	if err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair agendamento: %v", err))
	} else {
		logger.Info(fmt.Sprintf("✓ Agendamento: %s", agendamento))
	}
	
	logger.Info("✅ Busca concluída com sucesso!")
	return nil
}

// extractParticipantData - extrai todos os dados do participante
func (o *Orchestrator) extractParticipantData(ctx context.Context) (*models.ClientData, error) {
	logger.Info("========================================")
	logger.Info("ETAPA 3: EXTRAÇÃO DE DADOS DO PARTICIPANTE")
	logger.Info("========================================")
	
	// Clica em Participantes
	if err := o.participantsNav.ClickParticipantes(ctx, o.iframeWaiter); err != nil {
		return nil, err
	}
	
	// Extrai coobrigado
	cpfCoobrigado, nomeCoobrigado, err := o.participantsNav.ExtractCoobrigado(ctx, o.iframeWaiter)
	if err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair coobrigado: %v", err))
	}
	
	// Clica no CPF do proponente
	if err := o.participantsNav.ClickProponenteCPF(ctx, o.iframeWaiter); err != nil {
		return nil, err
	}
	
	// Aguarda página carregar
	time.Sleep(3 * time.Second)
	
	// Busca iframe dos detalhes
	iframeNode, err := o.iframeWaiter.WaitForIframe(ctx, "Detalhes Participante")
	if err != nil {
		return nil, err
	}
	
	// Extrai todos os dados do participante
	clientData, err := o.dataCoordinator.ExtractAllParticipantData(ctx, iframeNode)
	if err != nil {
		return nil, err
	}
	
	// Adiciona dados do coobrigado
	clientData.CoobrigadoCPF = cpfCoobrigado
	clientData.CoobrigadoNome = nomeCoobrigado
	
	logger.Info("✅ Dados do participante extraídos com sucesso!")
	return clientData, nil
}

// extractPropertyData - navega e extrai dados do imóvel
func (o *Orchestrator) extractPropertyData(ctx context.Context, clientData *models.ClientData) error {
	logger.Info("========================================")
	logger.Info("ETAPA 4: EXTRAÇÃO DE DADOS DO IMÓVEL")
	logger.Info("========================================")
	
	// Navega até página de imóvel
	if err := o.propertyNav.NavigateToProperty(ctx, o.menuNav, o.iframeWaiter); err != nil {
		return err
	}
	
	// Aguarda página carregar
	time.Sleep(3 * time.Second)
	
	// Busca iframe
	iframeNode, err := o.iframeWaiter.WaitForIframe(ctx, "Página Imóvel")
	if err != nil {
		return err
	}
	
	// Extrai dados do imóvel
	if err := o.dataCoordinator.ExtractPropertyData(ctx, iframeNode, clientData); err != nil {
		return err
	}
	
	logger.Info("✅ Dados do imóvel extraídos com sucesso!")
	return nil
}

// extractFinancialData - navega e extrai dados financeiros
func (o *Orchestrator) extractFinancialData(ctx context.Context, clientData *models.ClientData) error {
	logger.Info("========================================")
	logger.Info("ETAPA 5: EXTRAÇÃO DE DADOS FINANCEIROS")
	logger.Info("========================================")
	
	// Navega até valores da operação
	if err := o.propertyNav.NavigateToFinancialValues(ctx, o.menuNav, o.iframeWaiter); err != nil {
		return err
	}
	
	// Aguarda página carregar
	time.Sleep(3 * time.Second)
	
	// Busca iframe
	iframeNode, err := o.iframeWaiter.WaitForIframe(ctx, "Valores Operação")
	if err != nil {
		return err
	}
	
	// Extrai dados financeiros
	if err := o.dataCoordinator.ExtractFinancialData(ctx, iframeNode, clientData); err != nil {
		return err
	}
	
	logger.Info("✅ Dados financeiros extraídos com sucesso!")
	return nil
}

// logFinalResults - loga resultados finais
func (o *Orchestrator) logFinalResults(clientData *models.ClientData) {
	logger.Info("========================================")
	logger.Info("✅ PROCESSO CONCLUÍDO COM SUCESSO!")
	logger.Info("========================================")
	logger.Info(fmt.Sprintf("📝 Nome: %s", clientData.Nome))
	logger.Info(fmt.Sprintf("📋 CPF: %s", clientData.CPF))
	logger.Info(fmt.Sprintf("💼 Ocupação: %s", clientData.Ocupacao))
	logger.Info(fmt.Sprintf("🌍 Nacionalidade: %s", clientData.Nacionalidade))
	logger.Info(fmt.Sprintf("🆔 Tipo ID: %s | RG: %s", clientData.TipoIdentificacao, clientData.RG))
	logger.Info(fmt.Sprintf("👥 Coobrigado: %s (%s)", clientData.CoobrigadoNome, clientData.CoobrigadoCPF))
	logger.Info(fmt.Sprintf("📱 Telefone: %s", clientData.TelefoneCelular))
	logger.Info(fmt.Sprintf("🏠 Endereço: %s %s, %s - %s/%s (CEP: %s)", 
		clientData.TipoLogradouro, 
		clientData.Logradouro, 
		clientData.Numero, 
		clientData.Municipio, 
		clientData.UF,
		clientData.CEP))
	logger.Info(fmt.Sprintf("🏢 Endereço Imóvel: %s (CEP: %s)", clientData.EnderecoImovel, clientData.CEPImovel))
	logger.Info(fmt.Sprintf("💰 Valor Compra e Venda: %s", clientData.ValorCompraVenda))
	logger.Info(fmt.Sprintf("📄 Contrato: %s", clientData.NumeroContrato))
	logger.Info(fmt.Sprintf("💳 Conta: %s (Ag: %s)", clientData.ContaCorrente, clientData.Agencia))
	logger.Info("========================================")
}