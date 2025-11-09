package automation

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
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

// clickParticipanteCPF - clica no link do CPF do participante
func (bot *CaixaBot) clickParticipanteCPF(ctx context.Context) error {
	logger.Info("👤 Clicando no CPF do participante...")

	// Busca o iframe da página de Participantes
	iframeNode, err := bot.waitForIframe(ctx, "Página Participantes")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando link do CPF...")
			return nil
		}),

		chromedp.Sleep(2*time.Second),

		// Procura o link com onclick que contém "detalharParticipante"
		chromedp.WaitVisible(`a[onclick*="detalharParticipante"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Link do CPF encontrado!")
			return nil
		}),

		// Clica no link
		chromedp.Click(`a[onclick*="detalharParticipante"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ CPF clicado! Aguardando detalhes do participante...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}