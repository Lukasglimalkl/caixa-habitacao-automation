package config

import (
	"github.com/chromedp/chromedp"
)

// BrowserConfig - configurações do navegador
type BrowserConfig struct {
	Headless bool
	Options  []chromedp.ExecAllocatorOption
}

// DefaultBrowserConfig - configuração padrão do Chrome
func DefaultBrowserConfig(headless bool) BrowserConfig {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-software-rasterizer", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		chromedp.WindowSize(1920, 1080),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)
	
	return BrowserConfig{
		Headless: headless,
		Options:  opts,
	}
}