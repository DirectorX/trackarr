package processor

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"net/url"
	"strings"
)

/* Private */

func (p *Processor) processVarRule(node *xmlquery.Node, vars *map[string]string) error {
	result := ""

	// iterate elements in var node
	n := node.FirstChild
	for {
		// no more elements
		if n == nil {
			break
		}

		// skip junk nodes (mostly an empty line)
		nodeTag := strings.ToLower(strings.TrimSpace(n.Data))
		if nodeTag == "" {
			n = n.NextSibling
			continue
		}

		// process action
		switch nodeTag {
		case "string":
			// append value from element
			result += n.SelectAttr("value")
		case "var", "varenc":
			// append var
			varName := n.SelectAttr("name")
			varValue, ok := (*vars)[varName]
			if !ok {
				// check config
				varValue, ok = p.cfg.Config[varName]
			}

			if !ok {
				p.log.Errorf("Missing variable: %q", varName)
				return fmt.Errorf("missing variable: %q", varName)
			}

			// url encode value?
			if nodeTag == "varenc" {
				varValue = url.QueryEscape(varValue)
			}

			result += varValue

		default:
			p.log.Warnf("Unsupported var operation %q", nodeTag)
			break
		}

		// next element
		n = n.NextSibling
	}

	p.log.Debugf("Result for rule var: %s=%s", node.SelectAttr("name"), result)
	return nil
}
