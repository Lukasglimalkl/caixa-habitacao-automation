package extractors

import (
	"context"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/lukasglimalkl/caixa-habitacao-automation/rpa-service/pkg/logger"
)

// ExtractFieldFromTable - extrai valor de um campo da tabela usando o label
// Padrão da tabela: <tr><td><label>Campo:</label></td><td>Valor</td></tr>
func ExtractFieldFromTable(ctx context.Context, iframeNode *cdp.Node, labelText string) (string, error) {
	// XPath genérico: encontra <tr> que contém o label e pega o td seguinte
	xpath := fmt.Sprintf(`//tr[.//label[contains(., '%s')]]/td[@class='alinha_esquerda']`, labelText)
	
	var value string
	err := chromedp.Text(xpath, &value, chromedp.BySearch, chromedp.FromNode(iframeNode)).Do(ctx)
	
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(value), nil
}

// ExtractFieldWithFallback - tenta extrair campo, retorna string vazia se falhar
func ExtractFieldWithFallback(ctx context.Context, iframeNode *cdp.Node, labelText string, fieldName string) string {
	logger.Info(fmt.Sprintf("🔍 Extraindo %s...", fieldName))
	
	value, err := ExtractFieldFromTable(ctx, iframeNode, labelText)
	
	if err != nil || value == "" {
		logger.Info(fmt.Sprintf("⚠️ %s não encontrado", fieldName))
		return ""
	}
	
	logger.Info(fmt.Sprintf("✓ %s: %s", fieldName, value))
	return value
}

// ScrollToTable - faz scroll até uma tabela específica
func ScrollToTable(ctx context.Context, tableHeaderText string) error {
	logger.Info(fmt.Sprintf("📜 Fazendo scroll até tabela '%s'...", tableHeaderText))
	
	jsCode := fmt.Sprintf(`
		(function() {
			const xpath = "//th[contains(., '%s')]";
			const result = document.evaluate(xpath, document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null);
			const element = result.singleNodeValue;
			
			if (element) {
				element.scrollIntoView({behavior: 'smooth', block: 'center'});
				return true;
			}
			return false;
		})();
	`, tableHeaderText)
	
	var scrollSuccess bool
	err := chromedp.Evaluate(jsCode, &scrollSuccess).Do(ctx)
	
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao executar scroll: %v", err))
		return err
	}
	
	if scrollSuccess {
		logger.Info("✓ Scroll executado!")
	} else {
		logger.Info(fmt.Sprintf("⚠️ Tabela '%s' não encontrada", tableHeaderText))
	}
	
	return nil
}