package parser

import (
	"gitlab.com/cloudb0x/trackarr/config"

	"github.com/antchfx/xmlquery"
	"github.com/pkg/errors"
)

/* Private */

func parseRules(t *config.TrackerInfo, doc *xmlquery.Node) error {
	// this function is only responsible for grabbing the linematched xml node
	// the actual parsing / processing will happen via the processor package
	rules := xmlquery.FindOne(doc, "//parseinfo/linematched")
	if rules == nil {
		log.Errorf("Failed parsing tracker linematched rules")
		return errors.New("failed to parse tracker line matched rules")
	}

	// store for later use
	t.LineMatchedRules = rules
	return nil
}
