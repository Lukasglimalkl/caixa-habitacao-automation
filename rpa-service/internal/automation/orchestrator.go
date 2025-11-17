package automation

import (
	"context"
	"fmt"
	"time"

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

func (o *Orchestrator) Execute(ctx context.Context, username, password, cpf string) (*models.ClientData, error) {
	logger.Info("🚀 Iniciando processo de automação completo...")
	logger.Info("========================================")
	
	// ETAPA 1: LOGIN
	logger.Info("ETAPA 1: LOGIN")
	logger.Info("========================================")
	if err := o.executeLogin(ctx, username, password); err != nil {
		return nil, fmt.Errorf("erro no login: %w", err)
	}
	logger.Info("✅ Login realizado com sucesso!")
	
	// ETAPA 2: BUSCA POR CPF
	logger.Info("========================================")
	logger.Info("ETAPA 2: BUSCA POR CPF")
	logger.Info("========================================")
	if err := o.executeSearch(ctx, cpf); err != nil {
		return nil, fmt.Errorf("erro na busca: %w", err)
	}
	logger.Info("✅ Busca concluída com sucesso!")
	
	// Cria clientData vazio
	clientData := &models.ClientData{}
	
	// ETAPA 3: EXTRAÇÃO DE VALORES DA OPERAÇÃO (PRIMEIRO!)
	logger.Info("========================================")
	logger.Info("ETAPA 3: EXTRAÇÃO DE VALORES DA OPERAÇÃO")
	logger.Info("========================================")
	if err := o.extractFinancialData(ctx, clientData); err != nil {
		logger.Error("⚠️ Erro ao extrair dados financeiros: " + err.Error())
		// Não retorna erro, continua
	}
	
	// ETAPA 4: EXTRAÇÃO DE DADOS DO PARTICIPANTE
	logger.Info("========================================")
	logger.Info("ETAPA 4: EXTRAÇÃO DE DADOS DO PARTICIPANTE")
	logger.Info("========================================")
	participantData, err := o.extractParticipantData(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao extrair dados do participante: %w", err)
	}
	
	// Merge dos dados
	clientData = participantData
	
	// ETAPA 5: EXTRAÇÃO DE DADOS DO IMÓVEL
	logger.Info("========================================")
	logger.Info("ETAPA 5: EXTRAÇÃO DE DADOS DO IMÓVEL")
	logger.Info("========================================")
	if err := o.extractPropertyData(ctx, clientData); err != nil {
		logger.Error("⚠️ Erro ao extrair dados do imóvel: " + err.Error())
		// Não retorna erro, continua
	}
	
	logger.Info("========================================")
	logger.Info("✅ AUTOMAÇÃO CONCLUÍDA COM SUCESSO!")
	logger.Info("========================================")
	
	return clientData, nil
}

// executeLogin - executa o processo de login
func (o *Orchestrator) executeLogin(ctx context.Context, username, password string) error {
	if err := o.loginNav.Login(ctx, username, password); err != nil {
		return err
	}
	
	if err := o.loginNav.VerifyLoginSuccess(ctx); err != nil {
		return err
	}
	
	return nil
}

// executeSearch - executa a busca por CPF
func (o *Orchestrator) executeSearch(ctx context.Context, cpf string) error {
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
	
	return nil
}

// extractFinancialData - navega e extrai dados financeiros (PRIMEIRA ETAPA!)
func (o *Orchestrator) extractFinancialData(ctx context.Context, clientData *models.ClientData) error {
	logger.Info("💰 Extraindo valores da operação...")
	
	// Navega para Valores da Operação
	err := o.propertyNav.NavigateToFinancialValues(ctx, o.menuNav, o.iframeWaiter)
	if err != nil {
		return fmt.Errorf("erro ao navegar para valores da operação: %w", err)
	}
	
	logger.Info("✅ Navegou para Valores da Operação!")
	
	// Aguarda página carregar
	time.Sleep(3 * time.Second)
	
	// Aguarda iframe da página de valores
	iframeNode, err := o.iframeWaiter.WaitForIframe(ctx, "Valores da Operação")
	if err != nil {
		return fmt.Errorf("erro ao aguardar iframe: %w", err)
	}
	
	// Extrai dados financeiros
	err = o.dataCoordinator.ExtractFinancialData(ctx, iframeNode, clientData)
	if err != nil {
		return fmt.Errorf("erro ao extrair valores: %w", err)
	}
	
	logger.Info("✅ Valores da operação extraídos com sucesso!")
	return nil
}

// extractParticipantData - extrai todos os dados do participante
func (o *Orchestrator) extractParticipantData(ctx context.Context) (*models.ClientData, error) {
	// Clica em "Ir Para" para abrir menu
	if err := o.menuNav.ClickIrPara(ctx, o.iframeWaiter); err != nil {
		return nil, fmt.Errorf("erro ao abrir menu: %w", err)
	}
	
	// Clica em Participantes
	if err := o.menuNav.ClickMenuOption(ctx, o.iframeWaiter, "Participantes", "participantePI"); err != nil {
		return nil, err
	}
	
	// Aguarda página carregar
	time.Sleep(3 * time.Second)
	
	// Busca iframe da página de participantes
	iframeNode, err := o.iframeWaiter.WaitForIframe(ctx, "Participantes")
	if err != nil {
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
	iframeNode, err = o.iframeWaiter.WaitForIframe(ctx, "Detalhes Participante")
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