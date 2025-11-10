package automation

import (
	"context"
	"fmt"
	"strings"
	"time"

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

		// 🎯 IMPORTANTE: Pega o PRIMEIRO link (proponente, não coobrigado)
		// Se precisar do primeiro especificamente, use tr[@id='Item1']
		chromedp.WaitVisible(`tr[@id='Item1']//a[contains(@onclick, 'detalharParticipante')]`, chromedp.BySearch, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Link do CPF do PROPONENTE encontrado!")
			return nil
		}),

		// Clica no link do PROPONENTE (primeiro)
		chromedp.Click(`tr[@id='Item1']//a[contains(@onclick, 'detalharParticipante')]`, chromedp.BySearch, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ CPF do PROPONENTE clicado! Aguardando detalhes...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}