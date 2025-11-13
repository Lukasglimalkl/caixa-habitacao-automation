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

// SearchNavigator - interface para navegação de busca
type SearchNavigator interface {
	SearchByCPF(ctx context.Context, cpf string) error
	ClickFirstResult(ctx context.Context) error
	ExtractAgendamentoAssinatura(ctx context.Context, iframeWaiter IframeWaiter) (string, error)
}

// CaixaSearchNavigator - implementação para busca na Caixa
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
	
	// CRÍTICO: Aguarda iframe carregar ANTES de qualquer coisa
	logger.Info("⏳ Aguardando iframe carregar...")
	time.Sleep(5 * time.Second)
	
	// Busca o iframe (como slice)
	var iframeNodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		chromedp.Nodes(`iframe[id^="frameConteudo"]`, &iframeNodes, chromedp.BySearch, chromedp.AtLeast(0)),
	)
	
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao buscar iframe: %v", err))
		return fmt.Errorf("iframe não encontrado: %w", err)
	}
	
	if len(iframeNodes) == 0 {
		logger.Error("❌ Nenhum iframe encontrado")
		return fmt.Errorf("iframe não encontrado")
	}
	
	iframeNode := iframeNodes[0]
	logger.Info("✓ Iframe encontrado!")
	
	// Agora sim, interage DENTRO do iframe
	err = chromedp.Run(ctx,
		// Aguarda campo CPF estar visível DENTRO DO IFRAME
		chromedp.WaitVisible(`#cpfCnpj`, chromedp.ByID, chromedp.FromNode(iframeNode)),
		
		// Limpa o campo primeiro
		chromedp.Clear(`#cpfCnpj`, chromedp.ByID, chromedp.FromNode(iframeNode)),
		
		// Preenche CPF
		chromedp.SendKeys(`#cpfCnpj`, cpf, chromedp.ByID, chromedp.FromNode(iframeNode)),
		
		chromedp.Sleep(1*time.Second),
		
		// Clica no botão de pesquisar
		chromedp.Click(`#btConsultarProposta`, chromedp.ByID, chromedp.FromNode(iframeNode)),
		
		// Aguarda resultados aparecerem
		chromedp.Sleep(5*time.Second),
	)
	
	if err != nil {
		return err
	}
	
	logger.Info("✓ Busca iniciada, aguardando resultados...")
	return nil
}

// ClickFirstResult - clica no primeiro resultado da busca
func (nav *CaixaSearchNavigator) ClickFirstResult(ctx context.Context) error {
	logger.Info("🎯 Clicando no primeiro resultado...")
	
	// Busca o iframe novamente
	var iframeNodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.Nodes(`iframe[id^="frameConteudo"]`, &iframeNodes, chromedp.BySearch, chromedp.AtLeast(0)),
	)
	
	if err != nil || len(iframeNodes) == 0 {
		return fmt.Errorf("iframe não encontrado para clicar no resultado")
	}
	
	iframeNode := iframeNodes[0]
	
	// Aguarda tabela de resultados aparecer DENTRO DO IFRAME
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`table.tb_lista`, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		
		chromedp.Sleep(2*time.Second),
		
		// Clica no primeiro link de CPF
		chromedp.Click(`(//a[contains(@href, 'javascript:selecionarProposta')])[1]`, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		
		// Aguarda próxima página carregar
		chromedp.Sleep(5*time.Second),
	)
	
	if err != nil {
		return err
	}
	
	logger.Info("✓ Primeiro resultado clicado!")
	return nil
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
	
	agendamento = strings.TrimSpace(agendamento)
	logger.Info(fmt.Sprintf("✓ Agendamento: %s", agendamento))
	return agendamento, nil
}