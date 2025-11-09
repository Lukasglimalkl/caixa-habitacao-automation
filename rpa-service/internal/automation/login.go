package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// doLogin - executa o login no portal
func (bot *CaixaBot) doLogin(ctx context.Context, username, password string) error {
	logger.Info("🔐 Executando login...")

	return chromedp.Run(ctx,
		chromedp.Navigate(portalURL),
		chromedp.Sleep(3*time.Second),

		chromedp.WaitVisible(`#username`, chromedp.ByID),
		chromedp.SendKeys(`#username`, username, chromedp.ByID),

		chromedp.WaitVisible(`#password`, chromedp.ByID),
		chromedp.SendKeys(`#password`, password, chromedp.ByID),

		chromedp.WaitVisible(`.btn_login`, chromedp.ByQuery),
		chromedp.Click(`.btn_login`, chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("⏳ Aguardando login processar...")
			return nil
		}),

		chromedp.Sleep(5*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var currentURL string
			chromedp.Location(&currentURL).Do(ctx)
			logger.Info(fmt.Sprintf("✅ Login concluído - URL: %s", currentURL))
			return nil
		}),
	)
}