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

// clickParticipantes - clica no botão Participantes
func (bot *CaixaBot) clickParticipantes(ctx context.Context) error {
	logger.Info("👥 Clicando em Participantes...")

	// Usa o iframe da página de detalhes
	iframeNode, err := bot.waitForIframe(ctx, "Detalhes - Participantes")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando botão Participantes...")
			return nil
		}),

		// Espera o div aparecer
		chromedp.WaitVisible(`#participantePIDesabCheck`, chromedp.ByID, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão Participantes encontrado!")
			return nil
		}),

		// Clica no div
		chromedp.Click(`#participantePIDesabCheck`, chromedp.ByID, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Participantes clicado! Aguardando nova página...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}

// extractCoobrigadoFromTable - extrai CPF e Nome do coobrigado da TABELA (antes de clicar)
func (bot *CaixaBot) extractCoobrigadoFromTable(ctx context.Context, clientData *models.ClientData) error {
	logger.Info("📋 Extraindo dados do Coobrigado da tabela...")

	// Busca o iframe da página de Participantes
	iframeNode, err := bot.waitForIframe(ctx, "Tabela Participantes")
	if err != nil {
		return err
	}

	err = chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando linha do Coobrigado (Item2)...")

			// XPath para pegar o CPF do segundo participante (td com link)
			xpathCPF := `//tr[@id='Item2']//td[@class='alinha_centro'][1]//a[@class='link_normal']`
			
			var cpfCoobrigado string
			err := chromedp.Text(xpathCPF, &cpfCoobrigado, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("❌ Erro ao extrair CPF do coobrigado: %v", err))
				return err
			}

			clientData.CoobrigadoCPF = strings.TrimSpace(cpfCoobrigado)
			logger.Info(fmt.Sprintf("✓ CPF Coobrigado: %s", clientData.CoobrigadoCPF))

			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			// XPath para pegar o Nome do segundo participante (td seguinte)
			xpathNome := `//tr[@id='Item2']//td[@class='alinha_centro'][2]`
			
			var nomeCoobrigado string
			err := chromedp.Text(xpathNome, &nomeCoobrigado, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("❌ Erro ao extrair Nome do coobrigado: %v", err))
				return err
			}

			clientData.CoobrigadoNome = strings.TrimSpace(nomeCoobrigado)
			logger.Info(fmt.Sprintf("✓ Nome Coobrigado: %s", clientData.CoobrigadoNome))

			return nil
		}),
	)

	return err // ✅ AGORA ESTÁ CORRETO
}

// clickParticipanteCPF - clica no link do CPF do participante (PROPONENTE - primeiro da lista)
func (bot *CaixaBot) clickParticipanteCPF(ctx context.Context) error {
	logger.Info("👤 Clicando no CPF do PROPONENTE (primeiro participante)...")

	// Busca o iframe da página de Participantes
	iframeNode, err := bot.waitForIframe(ctx, "Página Participantes")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando link do CPF do PROPONENTE...")
			return nil
		}),

		chromedp.Sleep(2*time.Second),

		// 🎯 SELETOR MAIS GENÉRICO: pega o PRIMEIRO link com onclick="detalharParticipante"
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Buscando todos os links de participantes...")
			
			// Tenta encontrar o primeiro link
			var nodes []*cdp.Node
			err := chromedp.Nodes(`a[onclick*="detalharParticipante"]`, &nodes, chromedp.ByQueryAll, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("❌ Erro ao buscar links: %v", err))
				return err
			}
			
			logger.Info(fmt.Sprintf("✓ Encontrados %d participantes", len(nodes)))
			
			if len(nodes) == 0 {
				return fmt.Errorf("nenhum link de participante encontrado")
			}
			
			return nil
		}),

		// Espera o PRIMEIRO link aparecer (índice [0])
		chromedp.WaitVisible(`a[onclick*="detalharParticipante"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Link do primeiro participante encontrado!")
			return nil
		}),

		// Clica no PRIMEIRO link encontrado
		chromedp.Click(`a[onclick*="detalharParticipante"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ CPF do PROPONENTE clicado! Aguardando detalhes...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}

// clickParticipanteCPFWithRetry - tenta clicar no CPF com retry
func (bot *CaixaBot) clickParticipanteCPFWithRetry(ctx context.Context) error {
	maxRetries := 3
	
	for i := 0; i < maxRetries; i++ {
		logger.Info(fmt.Sprintf("👤 Tentativa %d/%d de clicar no CPF do PROPONENTE...", i+1, maxRetries))
		
		err := bot.clickParticipanteCPF(ctx)
		
		if err == nil {
			return nil
		}
		
		logger.Error(fmt.Sprintf("⚠️ Tentativa %d falhou: %v", i+1, err))
		
		if i < maxRetries-1 {
			logger.Info("⏳ Aguardando antes de tentar novamente...")
			time.Sleep(3 * time.Second)
		}
	}
	
	return fmt.Errorf("falhou após %d tentativas", maxRetries)
}

// debugListParticipantes - lista todos os participantes (DEBUG)
func (bot *CaixaBot) debugListParticipantes(ctx context.Context) error {
	logger.Info("🔍 [DEBUG] Listando todos os participantes...")

	iframeNode, err := bot.waitForIframe(ctx, "Debug Participantes")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Pega TODOS os links
			var nodes []*cdp.Node
			err := chromedp.Nodes(`a[onclick*="detalharParticipante"]`, &nodes, chromedp.ByQueryAll, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err != nil {
				logger.Error(fmt.Sprintf("Erro ao listar: %v", err))
				return err
			}
			
			logger.Info(fmt.Sprintf("Total de participantes: %d", len(nodes)))
			
			// Lista cada um
			for i, node := range nodes {
				var text string
				chromedp.Text(".", &text, chromedp.ByQuery, chromedp.FromNode(node)).Do(ctx)
				logger.Info(fmt.Sprintf("  [%d] CPF: %s", i, text))
			}
			
			return nil
		}),
	)
}