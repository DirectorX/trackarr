package processor

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	listutils "github.com/l3uddz/trackarr/utils/lists"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

/* Private */

func (p *Processor) processExtractTagsRule(node *xmlquery.Node, vars *map[string]string) error {
	srcVar := node.SelectAttr("srcvar")
	srcSplit := node.SelectAttr("split")

	// validate we have the required attributes so far
	if srcVar == "" {
		return errors.New("no srcvar specified")
	} else if srcSplit == "" {
		return errors.New("no split specified")
	}

	// retrieve srcVar
	existingValue, ok := (*vars)[srcVar]
	if !ok {
		return fmt.Errorf("srcvar did not exist: %q", srcVar)
	}

	// set values slice
	tagValues := make([]string, 0)
	for _, v := range strings.Split(existingValue, srcSplit) {
		tagValues = append(tagValues, strings.TrimSpace(v))
	}

	// set results map
	results := make(map[string]string, 0)

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
		case "setvarif":
			// parse node attributes
			varName := n.SelectAttr("varName")
			varRegex := n.SelectAttr("regex")
			varValue := n.SelectAttr("value")
			varNewValue := n.SelectAttr("newValue")

			// validate we have the minimum required attributes
			if varName == "" {
				return fmt.Errorf("failed parsing varName from: %s", n.OutputXML(true))
			}

			// is this a regex setvarif ?
			if varRegex != "" {
				// regex based setvarif
				foundMatch := false

				// compile regex for matching
				rxp, err := regexp.Compile(varRegex)
				if err != nil {
					return errors.Wrapf(err, "regex was invalid: %s", varRegex)
				}

				// iterate tag values looking for regex match
				for _, val := range tagValues {
					if rxp.MatchString(val) {
						// there was a match
						results[varName] = val
						foundMatch = true
					}
				}

				if !foundMatch {
					p.log.Tracef("No match found for %q regex: %s", varName, varRegex)
				}
			} else if varValue != "" {
				// value based setvarif
				if listutils.StringListContains(tagValues, varValue, false) {
					results[varName] = varNewValue
				} else {
					p.log.Tracef("No match found for value %q in %q", varValue, existingValue)
				}
			} else {
				// unsupported setvarif logic
				return fmt.Errorf("unsupported setvarif operation: %s", n.OutputXML(true))
			}

		default:
			return fmt.Errorf("unsupported extracttags operation: %q", nodeTag)
		}

		// next element
		n = n.NextSibling
	}

	// were results found?
	if len(results) < 1 {
		return nil
	}

	// set result in vars map
	for k, v := range results {
		(*vars)[k] = v
	}

	p.log.Debugf("Result for extracttags rule: %q = %#v", srcVar, results)
	return nil
}
