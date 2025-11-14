package navigation

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	logger.Info(fmt.Sprintf("🔍 Iniciando busca por CPF: %s", cpf))
	
	// PASSO 1: SEMPRE aguarda iframe PRIMEIRO
	logger.Info("📍 PASSO 1: Aguardando iframe carregar...")
	iframeWaiter := NewIframeWaiter(nav.maxRetries, nav.timeouts)
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Busca CPF")
	
	if err != nil {
		logger.Error("❌ Iframe não encontrado!")
		return fmt.Errorf("erro ao aguardar iframe: %w", err)
	}
	
	logger.Info("✅ Iframe encontrado! Iniciando busca...")
	
	// PASSO 2: Busca campo CPF DENTRO do iframe
	logger.Info("📍 PASSO 2: Procurando campo CPF dentro do iframe...")
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(`#cpfCnpj`, chromedp.ByID, chromedp.FromNode(iframeNode)),
	)
	
	if err != nil {
		logger.Error("❌ Campo CPF não encontrado dentro do iframe!")
		return fmt.Errorf("campo CPF não encontrado: %w", err)
	}
	
	logger.Info("✅ Campo CPF encontrado!")
	
	// PASSO 3: Preenche CPF
	logger.Info("📍 PASSO 3: Preenchendo CPF...")
	err = chromedp.Run(ctx,
		chromedp.Clear(`#cpfCnpj`, chromedp.ByID, chromedp.FromNode(iframeNode)),
		chromedp.SendKeys(`#cpfCnpj`, cpf, chromedp.ByID, chromedp.FromNode(iframeNode)),
	)
	
	if err != nil {
		logger.Error("❌ Erro ao preencher CPF!")
		return err
	}
	
	logger.Info("✅ CPF preenchido!")
	
	// PASSO 4: Clica no botão de buscar
	logger.Info("📍 PASSO 4: Clicando no botão de busca...")
	err = chromedp.Run(ctx,
	chromedp.Sleep(1*time.Second),
	
	// Clica no link com onclick
	chromedp.Click(`//a[@onclick="executaConsulta('cpfCnpjProposta');"]`, chromedp.BySearch, chromedp.FromNode(iframeNode)),
	
	chromedp.Sleep(3*time.Second),
)
	if err != nil {
		logger.Error("❌ Erro ao clicar no botão de busca!")
		return err
	}
	
	logger.Info("✅ Busca realizada com sucesso! Aguardando resultados...")
	
	// PASSO 5: Aguarda resultados
	time.Sleep(3 * time.Second)
	
	return nil
}

// ClickFirstResult - clica no primeiro resultado da busca
func (nav *CaixaSearchNavigator) ClickFirstResult(ctx context.Context) error {
	logger.Info("🎯 Clicando no primeiro resultado...")
	
	// PASSO 1: Busca iframe novamente (página pode ter recarregado)
	logger.Info("📍 PASSO 1: Buscando iframe dos resultados...")
	iframeWaiter := NewIframeWaiter(nav.maxRetries, nav.timeouts)
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Resultados")
	
	if err != nil {
		logger.Error("❌ Iframe dos resultados não encontrado!")
		return fmt.Errorf("iframe não encontrado: %w", err)
	}
	
	logger.Info("✅ Iframe dos resultados encontrado!")
	
	// PASSO 2: Aguarda tabela de resultados aparecer
	logger.Info("📍 PASSO 2: Aguardando tabela de resultados...")
	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		chromedp.WaitVisible(`table.tb_lista`, chromedp.BySearch, chromedp.FromNode(iframeNode)),
	)
	
	if err != nil {
		logger.Error("❌ Tabela de resultados não encontrada!")
		return fmt.Errorf("tabela de resultados não encontrada: %w", err)
	}
	
	logger.Info("✅ Tabela de resultados encontrada!")
	
	// PASSO 3: Clica no primeiro link (número da proposta)
	logger.Info("📍 PASSO 3: Clicando no primeiro resultado...")
	
	// XPath para o link com onclick="executa('localizarProposta.do..."
	xpath := `//table[contains(@class, 'tb_lista')]//a[contains(@onclick, "localizarProposta.do")]`
	
	err = chromedp.Run(ctx,
		chromedp.Sleep(1*time.Second),
		chromedp.WaitVisible(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Click(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Sleep(5*time.Second),
	)
	
	if err != nil {
		logger.Error("❌ Erro ao clicar no resultado!")
		return err
	}
	
	logger.Info("✅ Primeiro resultado clicado! Aguardando próxima página...")
	return nil
}

// ExtractAgendamentoAssinatura - extrai data de agendamento de assinatura
func (nav *CaixaSearchNavigator) ExtractAgendamentoAssinatura(ctx context.Context, iframeWaiter IframeWaiter) (string, error) {
	logger.Info("📅 Extraindo data de agendamento de assinatura...")
	
	// Aguarda iframe
	iframeNode, err := iframeWaiter.WaitForIframe(ctx, "Proposta Selecionada")
	if err != nil {
		logger.Error("❌ Iframe não encontrado!")
		return "", err
	}
	
	logger.Info("✅ Iframe encontrado! Procurando agendamento...")
	
	// XPath para data de agendamento
	xpath := `//tr[.//label[contains(., 'Agendamento da Assinatura:')]]/td[@class='alinha_esquerda']`
	
	var agendamento string
	
	// IMPORTANTE: Precisa estar dentro de chromedp.Run()
	err = chromedp.Run(ctx,
		chromedp.WaitVisible(xpath, chromedp.BySearch, chromedp.FromNode(iframeNode)),
		chromedp.Text(xpath, &agendamento, chromedp.BySearch, chromedp.FromNode(iframeNode)),
	)
	
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao extrair agendamento: %v", err))
		return "", err
	}
	
	agendamento = strings.TrimSpace(agendamento)
	logger.Info(fmt.Sprintf("✅ Agendamento: %s", agendamento))
	return agendamento, nil
}