package config

import "time"

// Timeouts - configurações de timeout do sistema
type Timeouts struct {
	BrowserContext   time.Duration
	PageLoad         time.Duration
	IframeWait       time.Duration
	ElementWait      time.Duration
	NetworkIdle      time.Duration
	ScrollWait       time.Duration
	AfterClick       time.Duration
	BetweenRetries   time.Duration
}

// DefaultTimeouts - timeouts padrão do sistema
func DefaultTimeouts() Timeouts {
	return Timeouts{
		BrowserContext:  5 * time.Minute,  // Tempo total máximo
		PageLoad:        60 * time.Second, // Tempo para página carregar
		IframeWait:      30 * time.Second, // Tempo para iframe aparecer
		ElementWait:     10 * time.Second, // Tempo para elemento aparecer
		NetworkIdle:     5 * time.Second,  // Tempo para rede ficar idle
		ScrollWait:      2 * time.Second,  // Tempo após scroll
		AfterClick:      3 * time.Second,  // Tempo após clique
		BetweenRetries:  5 * time.Second,  // Tempo entre tentativas
	}
}

// MaxRetries - configurações de tentativas
type MaxRetries struct {
	IframeSearch       int
	ElementClick       int
	DataExtraction     int
	PageNavigation     int
}

// DefaultMaxRetries - tentativas padrão
func DefaultMaxRetries() MaxRetries {
	return MaxRetries{
		IframeSearch:    30,
		ElementClick:    5,
		DataExtraction:  3,
		PageNavigation:  3,
	}
}