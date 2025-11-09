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
		chromedp.Flag("headless", isHeadless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel2 := chromedp.NewContext(allocCtx)
	ctx, cancel3 := context.WithTimeout(ctx, timeout)

	cancelAll := func() {
		cancel3()
		cancel2()
		cancel1()
	}

	return ctx, cancelAll
}