package processor

import (
	"fmt"
	"regexp"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func (p *Processor) processVarReplaceRule(node *xmlquery.Node, vars map[string]string) error {
	srcVar := node.SelectAttr("srcvar")
	targetVar := node.SelectAttr("name")
	regexVar := node.SelectAttr("regex")
	replaceVar := node.SelectAttr("replace")

	// validate we parsed all of the required variables (better ways of below, however, wanted to return relevant error)
	if srcVar == "" {
		return errors.New("srcvar had no value")
	} else if targetVar == "" {
		return errors.New("name had no value")
	} else if regexVar == "" {
		return errors.New("regex had no value")
	}

	// ensure srcVar exists in vars map
	existingValue, ok := vars[srcVar]
	if !ok {
		return fmt.Errorf("srcvar var did not exist: %q", srcVar)
	}

	// validate provided regex expression
	rxp, err := regexp.Compile(regexVar)
	if err != nil {
		return errors.Wrapf(err, "regex was invalid: %s", regexVar)
	}

	// do replace
	result := rxp.ReplaceAllString(existingValue, replaceVar)

	// set result in vars map
	vars[targetVar] = result

	p.Log.Tracef("Result for varreplace rule: %q = %s", targetVar, result)
	return nil
}
