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

	// Inicializa clientData aqui
	clientData := &models.ClientData{}

	// 6. Extrai dados do Coobrigado da tabela
	if err := bot.extractCoobrigadoFromTable(ctx, clientData); err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair coobrigado: %v", err))
	}

	// 7. Clica no CPF do PROPONENTE (COM RETRY)
	if err := bot.clickParticipanteCPFWithRetry(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao clicar no CPF: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar no CPF: %v", err),
		}, err
	}
	// 8. Extrai todos os dados do PROPONENTE (incluindo telefone e endereço)
	proponenteData, err := bot.extractDadosParticipante(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao extrair dados: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao extrair dados: %v", err),
		}, err
	}

	// 9. Mescla os dados
	clientData.CPF = proponenteData.CPF
	clientData.Nome = proponenteData.Nome
	clientData.NumeroContrato = proponenteData.NumeroContrato
	clientData.ContaDebitoCompleta = proponenteData.ContaDebitoCompleta
	clientData.Agencia = proponenteData.Agencia
	clientData.ContaCorrente = proponenteData.ContaCorrente
	clientData.AgendamentoAssinatura = agendamento
	clientData.TelefoneCelular = proponenteData.TelefoneCelular
	clientData.CEP = proponenteData.CEP
	clientData.TipoLogradouro = proponenteData.TipoLogradouro
	clientData.Logradouro = proponenteData.Logradouro
	clientData.Numero = proponenteData.Numero
	clientData.Bairro = proponenteData.Bairro
	clientData.Municipio = proponenteData.Municipio
	clientData.UF = proponenteData.UF
	clientData.Complemento = proponenteData.Complemento

	
		// 10. 🆕 Clica no botão "Ir para" (abre o menu)
if err := bot.clickIrPara(ctx); err != nil {
	logger.Error(fmt.Sprintf("❌ Erro ao clicar em 'Ir para': %v", err))
	return &models.SearchResponse{
		Success: false,
		Message: fmt.Sprintf("Erro ao clicar em 'Ir para': %v", err),
	}, err
}

// 11. 🆕 Clica no menu "Imóvel" (tenta pelo dialog primeiro, depois fallback)
logger.Info("🏠 Clicando no menu Imóvel...")
if err := bot.clickMenuImovel(ctx); err != nil {
	logger.Error(fmt.Sprintf("❌ Erro ao clicar no menu 'Imóvel' pelo dialog: %v", err))
	logger.Info("🔄 Tentando método alternativo (clicar diretamente)...")
	
	// FALLBACK: Tenta clicar diretamente no botão Imóvel
	if err := bot.clickImovelDirectly(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Método alternativo também falhou: %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar no menu 'Imóvel': %v", err),
		}, err
	}
	
	logger.Info("✓ Método alternativo funcionou!")
}

	// 12. 🆕 Extrai dados do Imóvel
	logger.Info("🏠 Extraindo dados do Imóvel...")
	imovelData, err := bot.extractDadosImovel(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair dados do imóvel: %v", err))
	} else {
		clientData.EnderecoImovel = imovelData.EnderecoImovel
		clientData.CEPImovel = imovelData.CEPImovel
		
		logger.Info(fmt.Sprintf("✓ Endereço Imóvel: %s", imovelData.EnderecoImovel))
		logger.Info(fmt.Sprintf("✓ CEP Imóvel: %s", imovelData.CEPImovel))
	}

	// 13. 🆕 Clica novamente no botão "Ir para" (para acessar Valores da Operação)
	if err := bot.clickIrPara(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao clicar em 'Ir para' (2ª vez): %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar em 'Ir para': %v", err),
		}, err
	}

	// 14. 🆕 Clica em "Valores da Operação"
	if err := bot.clickValoresOperacao(ctx); err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao clicar em 'Valores da Operação': %v", err))
		return &models.SearchResponse{
			Success: false,
			Message: fmt.Sprintf("Erro ao clicar em 'Valores da Operação': %v", err),
		}, err
	}

	// 15. 🆕 Extrai Valor de Compra e Venda
	valorCompraVenda, err := bot.extractValorCompraVenda(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("⚠️ Erro ao extrair valor de compra e venda: %v", err))
	} else {
		clientData.ValorCompraVenda = valorCompraVenda
		logger.Info(fmt.Sprintf("✓ Valor Compra e Venda: %s", valorCompraVenda))
	}

	logger.Info("========================================")
	logger.Info("✅ PROCESSO CONCLUÍDO!")
	logger.Info(fmt.Sprintf("📝 Nome: %s", clientData.Nome))
	logger.Info(fmt.Sprintf("📋 CPF: %s", clientData.CPF))
	logger.Info(fmt.Sprintf("💼 Ocupação: %s", clientData.Ocupacao))
	logger.Info(fmt.Sprintf("🌍 Nacionalidade: %s", clientData.Nacionalidade))
	logger.Info(fmt.Sprintf("🆔 Tipo ID: %s | RG: %s", clientData.TipoIdentificacao, clientData.RG))
	logger.Info(fmt.Sprintf("👥 Coobrigado: %s (%s)", clientData.CoobrigadoNome, clientData.CoobrigadoCPF))
	logger.Info(fmt.Sprintf("📱 Telefone: %s", clientData.TelefoneCelular))
	logger.Info(fmt.Sprintf("🏠 Endereço Residencial: %s %s, %s - %s/%s", clientData.TipoLogradouro, clientData.Logradouro, clientData.Numero, clientData.Municipio, clientData.UF))
	logger.Info(fmt.Sprintf("🏢 Endereço Imóvel: %s (CEP: %s)", clientData.EnderecoImovel, clientData.CEPImovel))
	logger.Info(fmt.Sprintf("💰 Valor Compra e Venda: %s", clientData.ValorCompraVenda))
	logger.Info(fmt.Sprintf("📄 Contrato: %s", clientData.NumeroContrato))
	logger.Info(fmt.Sprintf("💳 Conta: %s (Ag: %s)", clientData.ContaCorrente, clientData.Agencia))
	logger.Info(fmt.Sprintf("📅 Agendamento: %s", clientData.AgendamentoAssinatura))
	logger.Info("========================================")

	return &models.SearchResponse{
		Success: true,
		Message: "Dados extraídos com sucesso",
		Data:    clientData,
	}, nil
}