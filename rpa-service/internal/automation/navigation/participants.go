package navigation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/config"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// ParticipantsNavigator - interface para navegação de participantes
type ParticipantsNavigator interface {
	ClickParticipantes(ctx context.Context, iframeWaiter IframeWaiter) error
	ExtractCoobrigado(ctx context.Context, iframeWaiter IframeWaiter) (cpf, nome string, err error)
	ClickProponenteCPF(ctx context.Context, iframeWaiter IframeWaiter) error
}

// CaixaParticipantsNavigator - implementação para participantes
type CaixaParticipantsNavigator struct {
	timeouts   config.Timeouts
	maxRetries config.MaxRetries
}

// NewCaixaParticipantsNavigator - cria novo navegador de participantes
func NewCaixaParticipantsNavigator(timeouts config.Timeouts, maxRetries config.MaxRetries) *CaixaParticipantsNavigator {
	return &CaixaParticipantsNavigator{
		timeouts:   timeouts,
		maxRetries: maxRetries,
	}
}

// ClickParticipantes - clica no menu Participantes
func (nav *CaixaParticipantsNavigator) ClickParticipantes(ctx context.Context, iframeWaiter IframeWaiter) error {
	logger.Info("👥 Clicando no menu Participantes...")
	
	// Aguarda iframe
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Menu Principal")
	if err != nil {
		return err
	}
	
	// Lista de seletores possíveis para o botão Participantes
	selectors := []string{
		"#participantePIDesabCheck",
		"#participantePI",
		"#participantePICheck",
		"#participantePIDesab",
	}
	
	// Tenta cada seletor
	for _, selector := range selectors {
		var nodes []*cdp.Node
		err := chromedp.Nodes(selector, &nodes, chromedp.ByID, chromedp.FromNode(iframeNode)).Do(ctx)
		
		if err == nil && len(nodes) > 0 {
			logger.Info(fmt.Sprintf("✓ Botão Participantes encontrado: %s", selector))
			
			err = chromedp.Click(selector, chromedp.ByID, chromedp.FromNode(iframeNode)).Do(ctx)
			if err == nil {
				logger.Info("✓ Clique realizado com sucesso!")
				time.Sleep(nav.timeouts.AfterClick)
				return nil
			}
		}
	}
	
	return fmt.Errorf("botão Participantes não encontrado")
}

// ExtractCoobrigado - extrai CPF e nome do coobrigado
func (nav *CaixaParticipantsNavigator) ExtractCoobrigado(ctx context.Context, iframeWaiter IframeWaiter) (string, string, error) {
	logger.Info("👥 Extraindo dados do Coobrigado...")
	
	// Context com timeout específico (não crítico)
	coobrigadoCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	iframeNode, err := iframeWaiter.WaitForIframe(coobrigadoCtx, "Tabela Participantes")
	if err != nil {
		logger.Info("⚠️ Erro ao buscar iframe para coobrigado, continuando sem coobrigado")
		return "", "", nil
	}
	
	// Extrai CPF e Nome do coobrigado
	cpf, nome, err := nav.extractCoobrigadoFromTable(coobrigadoCtx, iframeNode)
	if err != nil {
		logger.Info("⚠️ Coobrigado não encontrado ou não existe")
		return "", "", nil
	}
	
	logger.Info(fmt.Sprintf("✓ Coobrigado: %s (%s)", nome, cpf))
	return cpf, nome, nil
}

// extractCoobrigadoFromTable - extrai coobrigado da tabela
func (nav *CaixaParticipantsNavigator) extractCoobrigadoFromTable(ctx context.Context, iframeNode *cdp.Node) (string, string, error) {
	logger.Info("🔍 Procurando linha do Coobrigado (Item2)...")
	
	// XPath para linha do coobrigado (Item2)
	xpathCPF := `//tr[@id='Item2']//td[contains(@onclick, 'exibirDetalhesParticipante')]`
	
	var cpf string
	err := chromedp.Text(xpathCPF, &cpf, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err != nil {
		return "", "", fmt.Errorf("CPF do coobrigado não encontrado: %w", err)
	}
	
	cpf = strings.TrimSpace(cpf)
	logger.Info(fmt.Sprintf("✓ CPF Coobrigado: %s", cpf))
	
	// XPath para nome do coobrigado
	xpathNome := `//tr[@id='Item2']//td[2]`
	
	var nome string
	err = chromedp.Text(xpathNome, &nome, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err != nil {
		return cpf, "", fmt.Errorf("nome do coobrigado não encontrado: %w", err)
	}
	
	nome = strings.TrimSpace(nome)
	logger.Info(fmt.Sprintf("✓ Nome Coobrigado: %s", nome))
	
	return cpf, nome, nil
}

// ClickProponenteCPF - clica no CPF do proponente
func (nav *CaixaParticipantsNavigator) ClickProponenteCPF(ctx context.Context, iframeWaiter IframeWaiter) error {
	logger.Info("👤 Clicando no CPF do PROPONENTE...")
	
	for tentativa := 1; tentativa <= nav.maxRetries.ElementClick; tentativa++ {
		logger.Info(fmt.Sprintf("👤 Tentativa %d/%d", tentativa, nav.maxRetries.ElementClick))
		
		err := nav.clickProponenteCPFSingleAttempt(ctx, iframeWaiter)
		if err == nil {
			logger.Info("✅ Clique no CPF bem-sucedido!")
			return nil
		}
		
		logger.Error(fmt.Sprintf("⚠️ Tentativa %d falhou: %v", tentativa, err))
		
		if tentativa < nav.maxRetries.ElementClick {
			logger.Info("⏳ Aguardando antes de tentar novamente...")
			time.Sleep(nav.timeouts.BetweenRetries)
		}
	}
	
	return fmt.Errorf("falhou após %d tentativas", nav.maxRetries.ElementClick)
}

// clickProponenteCPFSingleAttempt - uma tentativa de clicar no CPF
func (nav *CaixaParticipantsNavigator) clickProponenteCPFSingleAttempt(ctx context.Context, iframeWaiter IframeWaiter) error {
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Página Participantes")
	if err != nil {
		return err
	}
	
	// XPath para o CPF do primeiro participante (proponente)
	xpath := `//tr[@id='Item1']//td[contains(@onclick, 'exibirDetalhesParticipante')]`
	
	return chromedp.Run(ctx,
		chromedp.WaitVisible(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Click(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Sleep(nav.timeouts.PageLoad),
	)
}