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
		
		// 🆕 Extrai Telefone Celular
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.extractTelefoneCelular(ctx, iframeNode, &clientData)
		}),

		// Scroll até tabela de endereço
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.scrollToEnderecoTable(ctx)
		}),

		chromedp.Sleep(2*time.Second),
		
		// 🆕 Extrai Endereço completo
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.extractEndereco(ctx, iframeNode, &clientData)
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

// extractTelefoneCelular - extrai o telefone celular do participante
func (bot *CaixaBot) extractTelefoneCelular(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("📱 Extraindo Telefone Celular...")
	
	// XPath para encontrar o telefone celular
	xpath := `//tr[.//label[contains(., 'Telefone Celular:')]]/td[@class='alinha_esquerda'][1]`
	
	var telefoneCelular string
	err := chromedp.Text(xpath, &telefoneCelular, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao extrair telefone celular: %v", err))
		return err
	}
	
	clientData.TelefoneCelular = strings.TrimSpace(telefoneCelular)
	logger.Info(fmt.Sprintf("✓ Telefone Celular: %s", clientData.TelefoneCelular))
	
	return nil
}

// scrollToEnderecoTable - faz scroll até a tabela de Endereço
func (bot *CaixaBot) scrollToEnderecoTable(ctx context.Context) error {
	logger.Info("📜 Fazendo scroll até tabela 'Endereço'...")
	
	jsCode := `
		(function() {
			// Procura pelo th que contém "Endereço"
			const xpath = "//th[contains(., 'Endereço') and not(contains(., 'Correspondência'))]";
			const result = document.evaluate(xpath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null);
			const element = result.singleNodeValue;
			
			if (element) {
				console.log("Tabela Endereço encontrada, fazendo scroll...");
				element.scrollIntoView({behavior: 'smooth', block: 'center'});
				return true;
			} else {
				console.log("Tabela Endereço não encontrada!");
				return false;
			}
		})();
	`
	
	var scrollSuccess bool
	err := chromedp.Evaluate(jsCode, &scrollSuccess).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao executar scroll: %v", err))
	} else if scrollSuccess {
		logger.Info("✓ Scroll até tabela Endereço executado!")
	} else {
		logger.Info("⚠️ Tabela 'Endereço' não encontrada para scroll")
	}
	
	return nil
}

// extractEndereco - extrai todos os dados de endereço (PRIMEIRO endereço, não o de correspondência)
func (bot *CaixaBot) extractEndereco(ctx context.Context, iframeNode *cdp.Node, clientData *models.ClientData) error {
	logger.Info("🏠 Extraindo dados de Endereço...")
	
	// XPath base: pega a primeira tabela com "Endereço" no header (não a de correspondência)
	baseXPath := `//table[.//th[contains(text(), 'Endereço') and not(contains(text(), 'Correspondência'))]]`
	
	err := chromedp.Run(ctx,
		// CEP
		chromedp.ActionFunc(func(ctx context.Context) error {
			xpath := baseXPath + `//tr[.//label[contains(., 'CEP:')]]/td[@class='alinha_esquerda'][1]`
			var cep string
			err := chromedp.Text(xpath, &cep, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao extrair CEP: %v", err))
				return err
			}
			clientData.CEP = strings.TrimSpace(cep)
			logger.Info(fmt.Sprintf("✓ CEP: %s", clientData.CEP))
			return nil
		}),
		
		// Tipo de Logradouro
		chromedp.ActionFunc(func(ctx context.Context) error {
			xpath := baseXPath + `//tr[.//label[contains(., 'Tipo de Logradouro:')]]/td[@class='alinha_esquerda'][last()]`
			var tipoLogradouro string
			err := chromedp.Text(xpath, &tipoLogradouro, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao extrair Tipo de Logradouro: %v", err))
				return err
			}
			clientData.TipoLogradouro = strings.TrimSpace(tipoLogradouro)
			logger.Info(fmt.Sprintf("✓ Tipo de Logradouro: %s", clientData.TipoLogradouro))
			return nil
		}),
		
		// Logradouro
		chromedp.ActionFunc(func(ctx context.Context) error {
			xpath := baseXPath + `//tr[.//label[contains(., 'Logradouro:')]]/td[@class='alinha_esquerda'][1]`
			var logradouro string
			err := chromedp.Text(xpath, &logradouro, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao extrair Logradouro: %v", err))
				return err
			}
			clientData.Logradouro = strings.TrimSpace(logradouro)
			logger.Info(fmt.Sprintf("✓ Logradouro: %s", clientData.Logradouro))
			return nil
		}),
		
		// Número
		chromedp.ActionFunc(func(ctx context.Context) error {
			xpath := baseXPath + `//tr[.//label[contains(., 'Número:')]]/td[@class='alinha_esquerda'][last()]`
			var numero string
			err := chromedp.Text(xpath, &numero, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao extrair Número: %v", err))
				return err
			}
			clientData.Numero = strings.TrimSpace(numero)
			logger.Info(fmt.Sprintf("✓ Número: %s", clientData.Numero))
			return nil
		}),
		
		// Bairro
		chromedp.ActionFunc(func(ctx context.Context) error {
			xpath := baseXPath + `//tr[.//label[contains(., 'Bairro:')]]/td[@class='alinha_esquerda'][last()]`
			var bairro string
			err := chromedp.Text(xpath, &bairro, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao extrair Bairro: %v", err))
				return err
			}
			clientData.Bairro = strings.TrimSpace(bairro)
			logger.Info(fmt.Sprintf("✓ Bairro: %s", clientData.Bairro))
			return nil
		}),
		
		// Município - UF
		chromedp.ActionFunc(func(ctx context.Context) error {
			xpath := baseXPath + `//tr[.//label[contains(., 'Município - UF:')]]/td[@class='alinha_esquerda']`
			var municipioUF string
			err := chromedp.Text(xpath, &municipioUF, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao extrair Município-UF: %v", err))
				return err
			}
			
			// Separa Município e UF (formato: "SAO PAULO - SP")
			municipioUF = strings.TrimSpace(municipioUF)
			partes := strings.Split(municipioUF, "-")
			if len(partes) >= 2 {
				clientData.Municipio = strings.TrimSpace(partes[0])
				clientData.UF = strings.TrimSpace(partes[1])
			} else {
				clientData.Municipio = municipioUF
			}
			
			logger.Info(fmt.Sprintf("✓ Município: %s", clientData.Municipio))
			logger.Info(fmt.Sprintf("✓ UF: %s", clientData.UF))
			return nil
		}),
		
		// Complemento (pode estar vazio)
		chromedp.ActionFunc(func(ctx context.Context) error {
			xpath := baseXPath + `//tr[.//label[contains(., 'Complemento:')]]/td[@class='alinha_esquerda'][1]`
			var complemento string
			chromedp.Text(xpath, &complemento, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			clientData.Complemento = strings.TrimSpace(complemento)
			if clientData.Complemento != "" {
				logger.Info(fmt.Sprintf("✓ Complemento: %s", clientData.Complemento))
			}
			return nil
		}),
	)
	
	return err
}