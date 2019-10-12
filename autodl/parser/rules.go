package parser

import (
	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func parseTrackerRules(doc *xmlquery.Node, tracker *TrackerInfo) error {
	// this function is only responsible for grabbing the linematched xml node
	// the actual parsing / processing will happen via the processor package
	rules := xmlquery.FindOne(doc, "//parseinfo/linematched")
	if rules == nil {
		log.Errorf("Failed parsing tracker linematched rules")
		return errors.New("failed to parse tracker line matched rules")
	}

	// store for later use
	tracker.LineMatchedRules = rules
	return nil
}
