package navigation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/automation/config"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// SearchNavigator - interface para navegação de busca
type SearchNavigator interface {
	SearchByCPF(ctx context.Context, cpf string) error
	ClickFirstResult(ctx context.Context) error
}

// CaixaSearchNavigator - implementação para busca no portal da Caixa
type CaixaSearchNavigator struct {
	timeouts   config.Timeouts
	maxRetries config.MaxRetries
}

// NewCaixaSearchNavigator - cria novo navegador de busca
func NewCaixaSearchNavigator(timeouts config.Timeouts, maxRetries config.MaxRetries) *CaixaSearchNavigator {
	return &CaixaSearchNavigator{
		timeouts:   timeouts,
		maxRetries: maxRetries,
	}
}

// SearchByCPF - busca por CPF no portal
func (nav *CaixaSearchNavigator) SearchByCPF(ctx context.Context, cpf string) error {
	logger.Info(fmt.Sprintf("🔍 Buscando CPF: %s", cpf))
	
	return chromedp.Run(ctx,
		// Aguarda campo de busca aparecer
		chromedp.WaitVisible(`#cpfCnpj`, chromedp.ByID),
		
		// Limpa campo e preenche CPF
		chromedp.Clear(`#cpfCnpj`, chromedp.ByID),
		chromedp.SendKeys(`#cpfCnpj`, cpf, chromedp.ByID),
		
		// Clica no botão de pesquisar
		chromedp.Click(`#btConsultarProposta`, chromedp.ByID),
		
		// Aguarda resultado
		chromedp.Sleep(nav.timeouts.AfterClick),
	)
}

// ClickFirstResult - clica no primeiro resultado da busca
func (nav *CaixaSearchNavigator) ClickFirstResult(ctx context.Context) error {
	logger.Info("🎯 Clicando no primeiro resultado...")
	
	return chromedp.Run(ctx,
		// Aguarda tabela de resultados
		chromedp.WaitVisible(`table.tb_lista`, chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
		
		// Clica no primeiro link de CPF
		chromedp.ActionFunc(func(ctx context.Context) error {
			return nav.clickFirstCPFLink(ctx)
		}),
		
		// Aguarda página carregar
		chromedp.Sleep(nav.timeouts.PageLoad),
	)
}

// clickFirstCPFLink - clica no primeiro link de CPF da tabela
func (nav *CaixaSearchNavigator) clickFirstCPFLink(ctx context.Context) error {
	// XPath: primeiro link que contém "javascript:selecionarProposta"
	xpath := `(//a[contains(@href, 'javascript:selecionarProposta')])[1]`
	
	return chromedp.Run(ctx,
		chromedp.WaitVisible(xpath, chromedp.BySearch),
		chromedp.Click(xpath, chromedp.BySearch),
	)
}

// ExtractAgendamentoAssinatura - extrai data de agendamento de assinatura
func (nav *CaixaSearchNavigator) ExtractAgendamentoAssinatura(ctx context.Context, iframeWaiter IframeWaiter) (string, error) {
	logger.Info("📅 Extraindo data de agendamento de assinatura...")
	
	// Aguarda iframe
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Proposta Selecionada")
	if err != nil {
		return "", err
	}
	
	// XPath para data de agendamento
	xpath := `//tr[.//label[contains(., 'Agendamento da Assinatura:')]]/td[@class='alinha_esquerda']`
	
	var agendamento string
	err = chromedp.Text(xpath, &agendamento, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao extrair agendamento: %v", err))
		return "", err
	}
	
	logger.Info(fmt.Sprintf("✓ Agendamento: %s", agendamento))
	return agendamento, nil
}