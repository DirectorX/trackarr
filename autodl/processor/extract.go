package processor

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
	"regexp"
)

/* Private */

func (p *Processor) processExtractRule(node *xmlquery.Node, vars *map[string]string) error {
	srcVar := node.SelectAttr("srcvar")
	regexNode := node.SelectElement("/regex")

	// validate we parsed all of the required variables (better ways of below, however, wanted to return relevant error)
	if srcVar == "" {
		return errors.New("srcvar had no value")
	} else if regexNode == nil {
		return errors.New("regex element not found")
	}

	// set logic defaults
	isOptional := false
	if node.SelectAttr("optional") == "true" {
		isOptional = true
	}

	// ensure srcVar exists in vars map
	existingValue, ok := (*vars)[srcVar]
	if !ok {
		if !isOptional {
			// this was not an optional extract var...
			return fmt.Errorf("non-optional srcvar var did not exist: %q", srcVar)
		}

		// it was optional
		return nil
	}

	// validate provided regex expression
	regexVar := regexNode.SelectAttr("value")
	if regexVar == "" {
		return fmt.Errorf("regex element had no value: %s", regexNode.OutputXML(true))
	}

	rxp, err := regexp.Compile(regexVar)
	if err != nil {
		return errors.Wrapf(err, "regex was invalid: %s", regexVar)
	}

	// build list of vars this regex captures (via groups)
	regexVars := make([]string, 0)
	for _, v := range node.SelectElements("/vars/var") {
		varName := v.SelectAttr("name")
		if varName == "" {
			return fmt.Errorf("failed parsing var name from: %s", v.OutputXML(true))
		}
		regexVars = append(regexVars, varName)
	}

	// retrieve regex groups
	matches := rxp.FindStringSubmatch(existingValue)
	matchPos := 1
	matchCount := len(matches)
	if matchPos >= matchCount {
		if !isOptional {
			return fmt.Errorf("regex returned no matches: %s", regexVar)
		} else {
			return nil
		}
	}

	p.Log.Tracef("Extract regex %q matched with %d groups", regexVar, matchCount)

	// process captured vars
	results := map[string]string{}
	for _, varName := range regexVars {
		// this should not occur - but we must ensure we dont try and access an out of bounds index
		if matchPos > matchCount {
			p.Log.Warnf("Failed parsing extract regex var %q from match group %d", varName, matchPos)
			continue
		}

		// add var to result
		results[varName] = matches[matchPos]
	}

	// set result in vars map
	for varName, varValue := range results {
		(*vars)[varName] = varValue
	}

	p.Log.Tracef("Result for extract rule: %q = %+v", srcVar, results)
	return nil
}
