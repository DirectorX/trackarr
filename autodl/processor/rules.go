package processor

import (
	"strings"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func (p *Processor) processRules(rules *xmlquery.Node, vars *map[string]string) error {
	p.Log.Tracef("Processing linematched rules against: %+v", vars)

	// iterate rules node
	n := rules.FirstChild
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
		p.Log.Tracef("Processing linematched rule: %q", nodeTag)
		switch nodeTag {
		case "var":
			// concat var from other vars
			if err := p.processVarRule(n, vars); err != nil {
				return errors.Wrapf(err, "failed processing var rule: %s", n.OutputXML(true))
			}

		case "varreplace":
			// replace text in a var
			if err := p.processVarReplaceRule(n, vars); err != nil {
				return errors.Wrapf(err, "failed processing varreplace rule: %s", n.OutputXML(true))
			}

		case "extract":
			// create multiple vars from a single regex
			if err := p.processExtractRule(n, vars); err != nil {
				return errors.Wrapf(err, "failed processing extract rule: %s", n.OutputXML(true))
			}

		case "extracttags":
			// set a var if regex matches a tag in a var
			if err := p.processExtractTagsRule(n, vars); err != nil {
				return errors.Wrapf(err, "failed processing extracttags rule: %s", n.OutputXML(true))
			}

		case "extractone":
			// extract one var from a list of regexes
			if err := p.processExtractOneRule(n, vars); err != nil {
				return errors.Wrapf(err, "failed processing extractone rule: %s", n.OutputXML(true))
			}

		case "setregex":
			// set a var if a regex matches
			if err := p.processSetRegexRule(n, vars); err != nil {
				return errors.Wrapf(err, "failed processing setregex rule: %s", n.OutputXML(true))
			}

		case "if":
			// if statement
			if err := p.processIfRule(n, vars); err != nil {
				return errors.Wrapf(err, "failed processing if rule: %s", n.OutputXML(true))
			}

		default:
			p.Log.Tracef("Unsupported linematched rule: %q", nodeTag)
		}

		// process next
		n = n.NextSibling
	}

	p.Log.Tracef("Finished processing linematched rules")
	return nil
}
