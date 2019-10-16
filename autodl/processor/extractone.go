package processor

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"strings"
)

/* Private */

func (p *Processor) processExtractOneRule(node *xmlquery.Node, vars *map[string]string) error {
	// iterate elements in extractone node
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

		switch nodeTag {
		case "extract":
			// run extract rule
			if err := p.processExtractRule(n, vars); err == nil {
				// the extract rule was successful, lets return (as we only ever want to complete one successfully)
				return nil
			}
		default:
			p.Log.Tracef("unsupported extractone operation: %q", nodeTag)
		}

		// next element
		n = n.NextSibling
	}

	// if we are here, all extracts failed
	return fmt.Errorf("failed finding any matches for extractone rules")
}
