package automation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/internal/models"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// extractDadosParticipante - extrai todos os dados do participante
func (bot *CaixaBot) extractDadosParticipante(ctx context.Context) (*models.ClientData, error) {
	logger.Info("📊 Extraindo dados do participante...")

	// Busca o iframe da página de detalhes do participante
	iframeNode, err := bot.waitForIframe(ctx, "Detalhes Participante")
	if err != nil {
		return nil, err
	}

	var clientData models.ClientData

	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		// Extrai Número do Contrato
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.extractNumeroContrato(ctx, iframeNode, &clientData)
		}),

		// Extrai CPF
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.extractCPF(ctx, iframeNode, &clientData)
		}),

		// Extrai Nome
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.extractNome(ctx, iframeNode, &clientData)
		}),

		// Scroll até tabela de dados da conta
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.scrollToContaTable(ctx)
		}),

		chromedp.Sleep(2*time.Second),

		// Extrai Conta de Débito
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.extractContaDebito(ctx, iframeNode, &clientData)
		}),
	)

	if err != nil {
		return nil, err
	}

	return &clientData, nil
}

// extractNumeroContrato - extrai o número do contrato
func (bot *CaixaBot) extractNumeroContrato(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🔍 Extraindo Número do Contrato...")
	
	xpath := `//td[contains(., 'N° do Contrato:')]/following-sibling::td[@class='alinha_esquerda']`
	
	var numeroContrato string
	err := chromedp.Text(xpath, &numeroContrato, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao extrair contrato: %v", err))
		return err
	}
	
	clientData.NumeroContrato = strings.TrimSpace(numeroContrato)
	logger.Info(fmt.Sprintf("✓ Contrato: %s", clientData.NumeroContrato))
	
	return nil
}

// extractCPF - extrai o CPF
func (bot *CaixaBot) extractCPF(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🔍 Extraindo CPF...")
	
	xpath := `//td[contains(., 'CPF:')]/following-sibling::td`
	
	var cpf string
	err := chromedp.Text(xpath, &cpf, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao extrair CPF: %v", err))
		return err
	}
	
	clientData.CPF = strings.TrimSpace(cpf)
	logger.Info(fmt.Sprintf("✓ CPF: %s", clientData.CPF))
	
	return nil
}

// extractNome - extrai o nome
func (bot *CaixaBot) extractNome(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🔍 Extraindo Nome...")
	
	xpath := `//label[contains(., 'Nome:')]/ancestor::td/following-sibling::td[@class='alinha_esquerda']`
	
	var nome string
	err := chromedp.Text(xpath, &nome, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao extrair nome: %v", err))
		return err
	}
	
	clientData.Nome = strings.TrimSpace(nome)
	logger.Info(fmt.Sprintf("✓ Nome: %s", clientData.Nome))
	
	return nil
}

// scrollToContaTable - faz scroll até a tabela de dados da conta
func (bot *CaixaBot) scrollToContaTable(ctx context.Context) error {
	logger.Info("📜 Fazendo scroll até tabela 'Dados da Conta - Débito'...")
	
	jsCode := `
		(function() {
			// Procura pelo thead que contém "Dados da Conta - Débito"
			const xpath = "//th[contains(., 'Dados da Conta - Débito')]";
			const result = document.evaluate(xpath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null);
			const element = result.singleNodeValue;
			
			if (element) {
				console.log("Tabela encontrada, fazendo scroll...");
				element.scrollIntoView({behavior: 'smooth', block: 'center'});
				return true;
			} else {
				console.log("Tabela não encontrada!");
				return false;
			}
		})();
	`
	
	var scrollSuccess bool
	err := chromedp.Evaluate(jsCode, &scrollSuccess).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao executar scroll: %v", err))
	} else if scrollSuccess {
		logger.Info("✓ Scroll até tabela executado!")
	} else {
		logger.Info("⚠️ Tabela 'Dados da Conta' não encontrada para scroll")
	}
	
	return nil
}

// extractContaDebito - extrai a conta de débito
func (bot *CaixaBot) extractContaDebito(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🔍 Extraindo Conta de Débito com XPath específico...")
	
	// XPath específico para encontrar a conta de débito
	xpath := `//tr[@class='linha_azul'][.//label[contains(., 'Conta de Débito:')]]/td[@class='alinha_esquerda fonte_laranja']`
	
	var contaDebito string
	err := chromedp.Text(xpath, &contaDebito, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("❌ Erro ao extrair com XPath específico: %v", err))
		logger.Info("🔍 Tentando XPath alternativo...")
		
		// XPath alternativo
		xpathAlt := `//td[@class='alinha_esquerda fonte_laranja' and contains(text(), '-') and contains(text(), '0347')]`
		err = chromedp.Text(xpathAlt, &contaDebito, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
		
		if err != nil {
			logger.Error(fmt.Sprintf("❌ XPath alternativo também falhou: %v", err))
			return fmt.Errorf("conta de débito não encontrada")
		}
	}
	
	clientData.ContaDebitoCompleta = strings.TrimSpace(contaDebito)
	logger.Info(fmt.Sprintf("✓ Conta completa extraída: %s", clientData.ContaDebitoCompleta))
	
	if clientData.ContaDebitoCompleta != "" {
		clientData.Agencia, clientData.ContaCorrente = separarContaDebito(clientData.ContaDebitoCompleta)
	} else {
		logger.Error("❌ Conta de débito está vazia!")
	}
	
	return nil
}