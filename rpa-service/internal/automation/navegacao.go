package automation

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// clickIrPara - clica no botão "Ir para" (ícone no topo)
func (bot *CaixaBot) clickIrPara(ctx context.Context) error {
	logger.Info("🎯 Clicando no botão 'Ir para'...")

	// Busca o iframe da página atual
	iframeNode, err := bot.waitForIframe(ctx, "Botão Ir Para")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando botão 'Ir para'...")
			return nil
		}),

		chromedp.Sleep(1*time.Second),

		// Procura a imagem com onclick que contém "divFluxogramaProposta"
		chromedp.WaitVisible(`img[onclick*="divFluxogramaProposta"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão 'Ir para' encontrado!")
			return nil
		}),

		// Clica na imagem
		chromedp.Click(`img[onclick*="divFluxogramaProposta"]`, chromedp.ByQuery, chromedp.FromNode(iframeNode)),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão 'Ir para' clicado! Aguardando menu aparecer...")
			return nil
		}),

		chromedp.Sleep(3*time.Second),
	)
}

// clickImovelDirectly - tenta clicar diretamente no botão Imóvel (sem passar pelo Ir Para)
// Útil quando o dialog não abre ou está em um estado diferente
func (bot *CaixaBot) clickImovelDirectly(ctx context.Context) error {
	logger.Info("🏠 Tentando clicar DIRETAMENTE no botão Imóvel (fallback)...")

	// Busca o iframe da página atual
	iframeNode, err := bot.waitForIframe(ctx, "Botão Imóvel Direto")
	if err != nil {
		return err
	}

	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando botão Imóvel diretamente na página...")
			
			// Lista de IDs possíveis
			possibleIDs := []string{
				"imovelPIDesabCheck",
				"imovelPI",
				"imovelPICheck",
				"imovelPIDesab",
			}
			
			// Tenta cada ID dentro do iframe
			for _, id := range possibleIDs {
				var nodes []*cdp.Node
				err := chromedp.Nodes(`#`+id, &nodes, chromedp.ByID, chromedp.FromNode(iframeNode)).Do(ctx)
				
				if err == nil && len(nodes) > 0 {
					// Verifica se está visível
					var isVisible bool
					err = chromedp.Evaluate(fmt.Sprintf(`
						(function() {
							var frames = document.getElementsByTagName('iframe');
							for (var i = 0; i < frames.length; i++) {
								try {
									var el = frames[i].contentDocument.getElementById('%s');
									if (el) {
										var style = window.getComputedStyle(el);
										return style.display !== 'none' && style.visibility !== 'hidden';
									}
								} catch(e) {}
							}
							return false;
						})()
					`, id), &isVisible).Do(ctx)
					
					if err == nil && isVisible {
						logger.Info(fmt.Sprintf("✓ Botão Imóvel encontrado diretamente: #%s", id))
						return chromedp.Click(`#`+id, chromedp.ByID, chromedp.FromNode(iframeNode)).Do(ctx)
					}
				}
			}
			
			// XPath como última tentativa
			logger.Info("🔍 Tentando XPath para botão Imóvel...")
			xpathImovel := `//div[contains(@onclick, 'chamarImovel') and not(contains(@style, 'display: none'))]`
			
			var imovelNodes []*cdp.Node
			err = chromedp.Nodes(xpathImovel, &imovelNodes, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			
			if err == nil && len(imovelNodes) > 0 {
				logger.Info("✓ Botão Imóvel encontrado via XPath!")
				return chromedp.Click(xpathImovel, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
			}
			
			return fmt.Errorf("botão Imóvel não encontrado diretamente")
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Botão Imóvel clicado diretamente!")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}

// waitForMenuDialog - espera o dialog do menu aparecer
func (bot *CaixaBot) waitForMenuDialog(ctx context.Context) error {
	logger.Info("⏳ Aguardando dialog do menu aparecer...")
	
	return chromedp.Run(ctx,
		chromedp.Sleep(2*time.Second),
		
		chromedp.ActionFunc(func(ctx context.Context) error {
			// Aguarda o dialog aparecer
			for i := 0; i < 10; i++ {
				logger.Info(fmt.Sprintf("🔍 Tentativa %d/10 - Procurando dialog #divFluxogramaProposta...", i+1))
				
				// Verifica se o dialog existe E está visível
				var dialogVisible bool
				err := chromedp.Evaluate(`
					(function() {
						var dialog = document.getElementById('divFluxogramaProposta');
						if (!dialog) {
							console.log('Dialog não encontrado');
							return false;
						}
						var style = window.getComputedStyle(dialog);
						var isVisible = style.display !== 'none' && style.visibility !== 'hidden';
						console.log('Dialog encontrado! Display:', style.display, 'Visibility:', style.visibility);
						return isVisible;
					})()
				`, &dialogVisible).Do(ctx)
				
				if err == nil && dialogVisible {
					logger.Info(fmt.Sprintf("✓ Dialog encontrado e visível! (tentativa %d)", i+1))
					
					// Debug: mostra conteúdo do dialog
					var dialogContent string
					chromedp.Evaluate(`
						(function() {
							var dialog = document.getElementById('divFluxogramaProposta');
							return 'HTML Length: ' + dialog.innerHTML.length + ' chars, ChildNodes: ' + dialog.childNodes.length;
						})()
					`, &dialogContent).Do(ctx)
					logger.Info("📋 Conteúdo do dialog: " + dialogContent)
					
					return nil
				}
				
				logger.Info(fmt.Sprintf("⏳ Dialog ainda não está visível, aguardando... (%d/10)", i+1))
				time.Sleep(1 * time.Second)
			}
			
			logger.Error("❌ Dialog não encontrado após 10 tentativas")
			
			// Debug final: lista todos os elementos visíveis na página
			var allVisibleDivs string
			chromedp.Evaluate(`
				(function() {
					var divs = document.querySelectorAll('div[id*="Fluxograma"], div[id*="fluxograma"]');
					var result = [];
					for (var i = 0; i < divs.length; i++) {
						var style = window.getComputedStyle(divs[i]);
						result.push({
							id: divs[i].id,
							display: style.display,
							visibility: style.visibility
						});
					}
					return JSON.stringify(result, null, 2);
				})()
			`, &allVisibleDivs).Do(ctx)
			logger.Info("📋 Debug - Todos os divs com 'Fluxograma' no ID:")
			logger.Info(allVisibleDivs)
			
			return fmt.Errorf("dialog não encontrado após 10 tentativas")
		}),
	)
}

// clickMenuImovel - clica no menu "Imóvel" (dentro do dialog que pode ter iframe)
func (bot *CaixaBot) clickMenuImovel(ctx context.Context) error {
	logger.Info("🏠 Clicando no menu 'Imóvel'...")

	return chromedp.Run(ctx,
		// Aguarda o dialog aparecer
		chromedp.ActionFunc(func(ctx context.Context) error {
			return bot.waitForMenuDialog(ctx)
		}),
		
		chromedp.Sleep(2*time.Second),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("🔍 Procurando menu 'Imóvel' com múltiplas estratégias...")
			
			// Lista de IDs possíveis do botão Imóvel (do mais provável ao menos)
			possibleIDs := []string{
				"imovelPIDesabCheck",  // Desabilitado com Check (mais comum)
				"imovelPI",            // Habilitado normal
				"imovelPICheck",       // Habilitado com Check
				"imovelPIDesab",       // Desabilitado
				"imovelPICadeado",     // Com cadeado
			}
			
			// ESTRATÉGIA 1: Procura todos os IDs possíveis na página principal
			logger.Info("1️⃣ Procurando por todos os IDs possíveis do botão Imóvel...")
			for _, id := range possibleIDs {
				var nodes []*cdp.Node
				err := chromedp.Nodes(`#`+id, &nodes, chromedp.ByID).Do(ctx)
				
				if err == nil && len(nodes) > 0 {
					// Verifica se está visível (não tem display:none)
					var isVisible bool
					err = chromedp.Evaluate(fmt.Sprintf(`
						(function() {
							var el = document.getElementById('%s');
							if (!el) return false;
							var style = window.getComputedStyle(el);
							return style.display !== 'none' && style.visibility !== 'hidden';
						})()
					`, id), &isVisible).Do(ctx)
					
					if err == nil && isVisible {
						logger.Info(fmt.Sprintf("✓ Botão encontrado e visível: #%s", id))
						logger.Info(fmt.Sprintf("🎯 Clicando em #%s...", id))
						return chromedp.Click(`#`+id, chromedp.ByID).Do(ctx)
					}
					
					logger.Info(fmt.Sprintf("⚠️ Botão #%s existe mas não está visível", id))
				}
			}
			
			// ESTRATÉGIA 2: Procura iframe DENTRO do dialog
			logger.Info("2️⃣ Procurando iframe dentro do dialog...")
			var dialogIframes []*cdp.Node
			err := chromedp.Nodes(`#divFluxogramaProposta iframe`, &dialogIframes, chromedp.ByQueryAll).Do(ctx)
			
			if err == nil && len(dialogIframes) > 0 {
				logger.Info(fmt.Sprintf("✓ Iframe encontrado no dialog! Total: %d", len(dialogIframes)))
				dialogIframeNode := dialogIframes[0]
				
				// Tenta os mesmos IDs dentro do iframe
				for _, id := range possibleIDs {
					var iframeNodes []*cdp.Node
					err = chromedp.Nodes(`#`+id, &iframeNodes, chromedp.ByID, chromedp.FromNode(dialogIframeNode)).Do(ctx)
					
					if err == nil && len(iframeNodes) > 0 {
						logger.Info(fmt.Sprintf("✓ Botão encontrado no iframe: #%s", id))
						logger.Info(fmt.Sprintf("🎯 Clicando em #%s no iframe...", id))
						return chromedp.Click(`#`+id, chromedp.ByID, chromedp.FromNode(dialogIframeNode)).Do(ctx)
					}
				}
			}
			
			// ESTRATÉGIA 3: XPath mais agressivo (procura qualquer div com onclick="chamarImovel")
			logger.Info("3️⃣ Usando XPath para encontrar qualquer div com 'chamarImovel'...")
			xpathImovel := `//div[contains(@onclick, 'chamarImovel') and not(contains(@style, 'display: none'))]`
			
			var imovelNodes []*cdp.Node
			err = chromedp.Nodes(xpathImovel, &imovelNodes, chromedp.BySearch).Do(ctx)
			
			if err == nil && len(imovelNodes) > 0 {
				logger.Info(fmt.Sprintf("✓ Botão encontrado via XPath! Total encontrados: %d", len(imovelNodes)))
				logger.Info("🎯 Clicando no primeiro botão visível...")
				return chromedp.Click(xpathImovel, chromedp.BySearch).Do(ctx)
			}
			
			// ESTRATÉGIA 4: JavaScript direto (última tentativa)
			logger.Info("4️⃣ Tentando clicar via JavaScript direto...")
			var clicked bool
			err = chromedp.Evaluate(`
				(function() {
					var ids = ['imovelPIDesabCheck', 'imovelPI', 'imovelPICheck', 'imovelPIDesab'];
					for (var i = 0; i < ids.length; i++) {
						var el = document.getElementById(ids[i]);
						if (el && window.getComputedStyle(el).display !== 'none') {
							console.log('Clicando em: ' + ids[i]);
							el.click();
							return true;
						}
					}
					return false;
				})()
			`, &clicked).Do(ctx)
			
			if err == nil && clicked {
				logger.Info("✓ Clique executado via JavaScript!")
				return nil
			}
			
			// Se chegou aqui, nada funcionou
			logger.Error("❌ Menu Imóvel não encontrado em nenhuma das 4 estratégias!")
			
			// Debug: lista todos os divs visíveis
			var allDivs string
			chromedp.Evaluate(`
				(function() {
					var divs = document.querySelectorAll('div[id*="imovel"]');
					var result = [];
					for (var i = 0; i < divs.length; i++) {
						var style = window.getComputedStyle(divs[i]);
						result.push({
							id: divs[i].id,
							display: style.display,
							onclick: divs[i].getAttribute('onclick')
						});
					}
					return JSON.stringify(result, null, 2);
				})()
			`, &allDivs).Do(ctx)
			logger.Info("📋 Debug - Todos os divs com 'imovel' no ID:")
			logger.Info(allDivs)
			
			return fmt.Errorf("menu imóvel não encontrado após 4 estratégias")
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			logger.Info("✓ Menu 'Imóvel' clicado! Aguardando página carregar...")
			return nil
		}),

		chromedp.Sleep(4*time.Second),
	)
}