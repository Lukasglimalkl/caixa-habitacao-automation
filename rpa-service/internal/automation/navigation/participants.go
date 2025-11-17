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
	
	// Busca iframe
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Participantes")
	if err != nil {
		logger.Error("❌ Iframe não encontrado!")
		return err
	}
	
	logger.Info("✅ Iframe encontrado! Procurando botão...")
	
	// Lista de seletores possíveis para o botão Participantes
	selectors := []string{
		"#participantePIDesabCheck",
		"#participantePI",
		"#participantePICheck",
		"#participantePIDesab",
	}
	
	// Tenta cada seletor diretamente
	for _, selector := range selectors {
		logger.Info(fmt.Sprintf("🔍 Tentando seletor: %s", selector))
		
		// Tenta clicar direto
		err := chromedp.Run(ctx,
			chromedp.Sleep(1*time.Second),
			chromedp.Click(selector, chromedp.ByID, chromedp.FromNode(iframeNode)),
		)
		
		// Se conseguiu clicar, sucesso!
		if err == nil {
			logger.Info(fmt.Sprintf("✅ Botão Participantes clicado: %s", selector))
			time.Sleep(nav.timeouts.AfterClick)
			return nil
		}
		
		// Se falhou, tenta próximo
		logger.Info(fmt.Sprintf("⚠️ Seletor %s não funcionou: %v", selector, err))
	}
	
	logger.Error("❌ Botão Participantes não encontrado com nenhum seletor!")
	return fmt.Errorf("botão Participantes não encontrado")
}


// ExtractCoobrigado - extrai CPF e nome do coobrigado
func (nav *CaixaParticipantsNavigator) ExtractCoobrigado(ctx context.Context, iframeWaiter IframeWaiter) (string, string, error) {
	logger.Info("👥 Extraindo dados do Coobrigado...")
	
	// Context com timeout específico (não crítico)
	coobrigadoCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	// Busca iframe
	iframeNode, err := iframeWaiter.WaitForIframe(coobrigadoCtx, "Coobrigado")
	if err != nil {
		logger.Error("❌ Iframe não encontrado!")
		return "", "", err
	}
	
	logger.Info("✅ Iframe encontrado! Procurando coobrigado...")
	
	// Extrai CPF e Nome do coobrigado
	cpf, nome, err := nav.extractCoobrigadoFromTable(coobrigadoCtx, iframeNode)
	if err != nil {
		logger.Info("⚠️ Coobrigado não encontrado ou não existe")
		return "", "", nil
	}
	
	logger.Info(fmt.Sprintf("✅ Coobrigado encontrado: %s (%s)", nome, cpf))
	return cpf, nome, nil
}

// extractCoobrigadoFromTable - extrai coobrigado da tabela
func (nav *CaixaParticipantsNavigator) extractCoobrigadoFromTable(ctx context.Context, iframeNode *cdp.Node) (string, string, error) {
	logger.Info("🔍 Procurando linha do Coobrigado (Item2)...")
	
	// XPath CORRETO para o link com CPF do coobrigado
	xpathCPF := `//tr[@id='Item2']//a[contains(@onclick, 'detalharParticipante')]`
	
	var cpf string
	err := chromedp.Run(ctx,
		chromedp.WaitVisible(xpathCPF, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Text(xpathCPF, &cpf, chromedp.BySearch, chromedp.FromNode(iframeNode)),
	)
	
	if err != nil {
		return "", "", fmt.Errorf("CPF do coobrigado não encontrado: %w", err)
	}
	
	cpf = strings.TrimSpace(cpf)
	logger.Info(fmt.Sprintf("✅ CPF Coobrigado: %s", cpf))
	
	// XPath CORRETO para nome do coobrigado (coluna 3)
	xpathNome := `//tr[@id='Item2']//td[3]`
	
	var nome string
	err = chromedp.Run(ctx,
		chromedp.Text(xpathNome, &nome, chromedp.BySearch, chromedp.FromNode(iframeNode)),
	)
	
	if err != nil {
		return cpf, "", fmt.Errorf("nome do coobrigado não encontrado: %w", err)
	}
	
	nome = strings.TrimSpace(nome)
	logger.Info(fmt.Sprintf("✅ Nome Coobrigado: %s", nome))
	
	return cpf, nome, nil
}

// ClickProponenteCPF - clica no CPF do proponente
func (nav *CaixaParticipantsNavigator) ClickProponenteCPF(ctx context.Context, iframeWaiter IframeWaiter) error {
	logger.Info("👤 Clicando no CPF do PROPONENTE...")
	
	// Busca iframe UMA VEZ SÓ
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Proponente")
	if err != nil {
		logger.Error("❌ Iframe não encontrado!")
		return err
	}
	
	logger.Info("✅ Iframe encontrado! Procurando CPF do proponente...")
	
	// XPath para o link com CPF do proponente
	xpath := `//tr[@id='Item1']//a[contains(@onclick, 'detalharParticipante')]`
	
	// Tenta clicar com retries
	for tentativa := 1; tentativa <= nav.maxRetries.ElementClick; tentativa++ {
		logger.Info(fmt.Sprintf("🎯 Tentativa %d/%d de clicar no CPF", tentativa, nav.maxRetries.ElementClick))
		
		err := chromedp.Run(ctx,
			chromedp.WaitVisible(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
			chromedp.Click(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
			chromedp.Sleep(3*time.Second),
		)
		
		if err == nil {
			logger.Info("✅ Clique realizado! Verificando se página carregou...")
			
			// Verifica se entrou na página de detalhes
			if nav.verifyDetailPageLoaded(ctx, iframeWaiter) {
				logger.Info("✅ Página 'Detalhe do Participante' carregada com sucesso!")
				return nil
			}
			
			logger.Info("⚠️ Página de detalhes não carregou, tentando novamente...")
		} else {
			logger.Error(fmt.Sprintf("⚠️ Tentativa %d falhou ao clicar: %v", tentativa, err))
		}
		
		if tentativa < nav.maxRetries.ElementClick {
			logger.Info("⏳ Aguardando antes de tentar novamente...")
			time.Sleep(nav.timeouts.BetweenRetries)
		}
	}
	
	logger.Error(fmt.Sprintf("❌ Todas as %d tentativas falharam!", nav.maxRetries.ElementClick))
	return fmt.Errorf("falhou após %d tentativas", nav.maxRetries.ElementClick)
}

// verifyDetailPageLoaded - verifica se a página de detalhes carregou
func (nav *CaixaParticipantsNavigator) verifyDetailPageLoaded(ctx context.Context, iframeWaiter IframeWaiter) bool {
	logger.Info("🔍 Verificando se página de detalhes carregou...")
	
	// Aguarda um pouco para a página carregar
	time.Sleep(2 * time.Second)
	
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Detalhe Participante")
	if err != nil {
		logger.Error("❌ Iframe de detalhes não encontrado")
		return false
	}
	
	// Procura pelo título específico da página
	xpath := `//h1//span[@class='subtitulo_paginas' and contains(., 'Detalhe do Participante')]`
	
	var nodes []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.Nodes(xpath, &nodes, chromedp.BySearch, chromedp.FromNode(iframeNode)),
	)
	
	if err == nil && len(nodes) > 0 {
		logger.Info("✅ Título 'Detalhe do Participante' encontrado!")
		return true
	}
	
	logger.Error("❌ Título 'Detalhe do Participante' não encontrado")
	return false
}