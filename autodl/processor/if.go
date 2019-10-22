package processor

import (
	"fmt"
	"regexp"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func (p *Processor) processIfRule(node *xmlquery.Node, vars map[string]string) error {
	srcVar := node.SelectAttr("srcvar")
	varRegex := node.SelectAttr("regex")

	// validate we have a new var name
	if srcVar == "" {
		return errors.New("no srcvar specified")
	} else if varRegex == "" {
		return errors.New("no regex specified")
	}

	// retrieve srcVar
	existingValue, ok := vars[srcVar]
	if !ok {
		return fmt.Errorf("srcvar did not exist: %q", srcVar)
	}

	// compile regex for matching
	rxp, err := regexp.Compile(varRegex)
	if err != nil {
		return errors.Wrapf(err, "regex was invalid: %s", varRegex)
	}

	// if condition met?
	if rxp.MatchString(existingValue) {
		// condition was met, process it's ruleset
		return p.processRules(node, vars)
	}

	return nil
}
