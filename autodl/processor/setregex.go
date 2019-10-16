package processor

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	"regexp"
)

/* Private */

func (p *Processor) processSetRegexRule(node *xmlquery.Node, vars *map[string]string) error {
	srcVar := node.SelectAttr("srcvar")
	regexVar := node.SelectAttr("regex")
	targetVar := node.SelectAttr("varName")
	targetVal := node.SelectAttr("newValue")

	// validate we parsed all of the required variables (better ways of below, however, wanted to return relevant error)
	if srcVar == "" {
		return errors.New("srcvar had no value")
	} else if targetVar == "" {
		return errors.New("varName had no value")
	} else if regexVar == "" {
		return errors.New("regex had no value")
	}

	// ensure srcVar exists in vars map
	existingValue, ok := (*vars)[srcVar]
	if !ok {
		return fmt.Errorf("srcvar var did not exist: %q", srcVar)
	}

	// validate provided regex expression
	rxp, err := regexp.Compile(regexVar)
	if err != nil {
		return errors.Wrapf(err, "regex was invalid: %s", regexVar)
	}

	// does regex match?
	if !rxp.MatchString(existingValue) {
		return nil
	}

	// set result in vars map
	(*vars)[targetVar] = targetVal

	p.Log.Tracef("Result for setregex rule: %q = %s", targetVar, targetVal)
	return nil
}
