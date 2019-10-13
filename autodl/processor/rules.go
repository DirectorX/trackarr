package processor

import (
	"strings"
)

/* Private */

func (p *Processor) processRules(vars *map[string]string) error {
	p.log.Tracef("Processing linematched rules against: %#v", vars)

	// iterate linematched (rules) node
	n := p.tracker.LineMatchedRules.FirstChild
	for {
		// break when node is empty
		if n == nil {
			break
		}

		// skip junk nodes (mostly an empty line)
		nodeTag := strings.ToLower(strings.TrimSpace(n.Data))
		if nodeTag == "" {
			n = n.NextSibling
			continue
		}

		// process tag
		p.log.Tracef("Processing linematched rule: %q", nodeTag)
		switch nodeTag {
		case "var":
			// concat var from other vars
			if err := p.processVarRule(n, vars); err != nil {
				p.log.WithError(err).Error("failed processing var rule")
				return err
			}

		case "varreplace":
			// replace text in a var
			break
		case "extract":
			// create multiple vars from a single regex
			break
		case "extracttags":
			// set a var if regex matches a tag in a var
			break
		case "extractone":
			// extract one var from a list of regexes
			break
		case "setregex":
			// set a var if a regex matches
			break
		case "if":
			// if statement
			break
		default:
			p.log.Warnf("Unsupported linematched rule: %q", nodeTag)
		}

		// process next
		n = n.NextSibling
	}

	p.log.Tracef("Finished processing linematched rules")
	return nil
}
