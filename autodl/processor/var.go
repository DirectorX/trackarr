package processor

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func (p *Processor) processVarRule(node *xmlquery.Node, vars map[string]string) error {
	result := ""
	newVarName := node.SelectAttr("name")

	// validate we have a new var name
	if newVarName == "" {
		return errors.New("no new var name specified")
	}

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
			// append value to result
			result += n.SelectAttr("value")
		case "var", "varenc":
			// append var
			varName := n.SelectAttr("name")

			// did we have a var name to lookup?
			if varName == "" {
				return errors.New("var had no name to lookup")
			}

			// lookup var
			varValue, ok := vars[varName]
			if !ok {
				// do we have the variable in the user defined tracker config? (torrent_pass, passkey etc...)
				varValue, ok = p.Tracker.Config.Settings[varName]
				if !ok {
					return fmt.Errorf("missing variable: %q in %+v", varName, p.Tracker.Config.Settings)
				}
			}

			// url encode value?
			if nodeTag == "varenc" {
				varValue = url.QueryEscape(varValue)
			}

			// append value to result
			result += varValue

		default:
			p.Log.Tracef("unsupported var operation: %q", nodeTag)
		}

		// next element
		n = n.NextSibling
	}

	// set result in vars map
	vars[newVarName] = result

	p.Log.Tracef("Result for var rule: %q = %s", newVarName, result)
	return nil
}
