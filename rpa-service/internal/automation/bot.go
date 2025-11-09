package automation

import (
	"fmt"

	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

const (
	portalURL = "https://habitacao.caixa.gov.br/siopiweb-web/"
)

// CaixaBot - estrutura principal do bot
type CaixaBot struct{}

// NewCaixaBot - cria uma nova instância do bot
func NewCaixaBot() *CaixaBot {
	return &CaixaBot{}
}

// LoginAndSearch - função principal que orquestra todo o processo
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

	// 2. Busca CPF
	if err := bot.fillAndSearchCPF(ctx, req.CPF); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao buscar CPF: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao buscar CPF: %v", err),
		}, err
	}

	// 3. Clica na proposta
	if err := bot.clickProposta(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao clicar na proposta: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar na proposta: %v", err),
		}, err
	}

	// 4. Extrai agendamento
	agendamento, err := bot.extractAgendamento(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair agendamento: %v", err))
		agendamento = "Não encontrado"
	}

	// 5. Clica em Participantes
	if err := bot.clickParticipantes(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao clicar em Participantes: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar em Participantes: %v", err),
		}, err
	}

	// 6. Clica no CPF do participante
	if err := bot.clickParticipanteCPF(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao clicar no CPF: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar no CPF: %v", err),
		}, err
	}

	// 7. Extrai todos os dados do participante
	clientData, err := bot.extractDadosParticipante(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao extrair dados: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao extrair dados: %v", err),
		}, err
	}

	// 8. Adiciona o agendamento
	clientData.AgendamentoAssinatura = agendamento

	logger.Info("========================================")
	logger.Info("✅ PROCESSO CONCLUÍDO!")
	logger.Info(fmt.Sprintf("📝 Nome: %s", clientData.Nome))
	logger.Info(fmt.Sprintf("📋 CPF: %s", clientData.CPF))
	logger.Info(fmt.Sprintf("📄 Contrato: %s", clientData.NumeroContrato))
	logger.Info(fmt.Sprintf("🏦 Conta completa: %s", clientData.ContaDebitoCompleta))
	logger.Info(fmt.Sprintf("🏢 Agência: %s", clientData.Agencia))
	logger.Info(fmt.Sprintf("💳 Conta Corrente: %s", clientData.ContaCorrente))
	logger.Info(fmt.Sprintf("📅 Agendamento: %s", clientData.AgendamentoAssinatura))
	logger.Info("========================================")

	return &models.SearchResponse{
		Success: true,
		Message: "Dados extraídos com sucesso",
		Data:    clientData,
	}, nil
}