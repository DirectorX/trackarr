package processor

import (
	stringutils "github.com/l3uddz/trackarr/utils/strings"
	"strings"
)

/* Private */

func (p *Processor) processRules(vars *map[string]string) error {
	p.log.Tracef("Processing linematched rules against; %s", stringutils.JsonifyLax(vars))

	// iterate linematched (rules) node
	n := p.tracker.LineMatchedRules.FirstChild
	for {
		// break when node is empty
		if n == nil {
			break
		}

		// skip junk nodes (mostly an empty line)
		nodeTag := strings.TrimSpace(n.Data)
		if nodeTag == "" {
			n = n.NextSibling
			continue
		}

		// process tag
		p.log.Tracef("Processing linematched rule: %q", nodeTag)
		switch strings.ToLower(nodeTag) {
		case "var":
			// concat var from other vars
			break
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
