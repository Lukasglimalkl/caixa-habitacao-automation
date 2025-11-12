package automation

import (
	"context"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

const (
	timeout = 120 * time.Second
)

// createBrowserContext - cria contexto do navegador (reutilizável)
func (bot *CaixaBot) createBrowserContext() (context.Context, context.CancelFunc) {
	isHeadless := os.Getenv("HEADLESS") != "false"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// Configurações básicas
		chromedp.Flag("headless", isHeadless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		
		// 🆕 Flags antidetecção para headless funcionar melhor
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-web-security", false),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-infobars", true),
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("disable-background-networking", false),
		chromedp.Flag("enable-features", "NetworkService,NetworkServiceInProcess"),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("force-color-profile", "srgb"),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("safebrowsing-disable-auto-update", true),
		chromedp.Flag("password-store", "basic"),
		chromedp.Flag("use-mock-keychain", true),
		
		// User agent realista
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel2 := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(s string, i ...interface{}) {}))
	ctx, cancel3 := context.WithTimeout(ctx, timeout)

	cancelAll := func() {
		cancel3()
		cancel2()
		cancel1()
	}

	return ctx, cancelAll
}