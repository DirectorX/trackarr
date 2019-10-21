package release

import (
	"github.com/l3uddz/trackarr/config"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
)

/* Privates */
func (r *Release) shouldAccept(pvr *config.PvrInstance) (bool, error) {
	// iterate accept expressions
	for _, expression := range pvr.AcceptsExpr {
		result, err := expr.Run(expression, r.Info)
		if err != nil {
			return false, errors.Wrapf(err, "failed checking accept expression")
		}

		expResult, ok := result.(bool)
		if !ok {
			return false, errors.New("failed type asserting accept expression result")
		}

		if expResult {
			r.Log.Tracef("Allowing release for pvr %q due to accept expression match", pvr.Config.Name)
			return true, nil
		}
	}

	return false, nil
}
