package release

import (
	"gitlab.com/cloudb0x/trackarr/config"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
)

/* Privates */
func (r *Release) shouldIgnore(pvr *config.PvrInstance) (bool, error) {
	// iterate ignore expressions
	for _, expression := range pvr.IgnoresExpr {
		result, err := expr.Run(expression, r.Info)
		if err != nil {
			return true, errors.Wrap(err, "failed checking ignore expression")
		}

		expResult, ok := result.(bool)
		if !ok {
			return true, errors.New("failed type asserting ignore expression result")
		}

		if expResult {
			r.Log.Tracef("Ignoring release for pvr %q due to ignore expression match", pvr.Config.Name)
			return true, nil
		}
	}

	return false, nil
}
