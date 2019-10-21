package release

import (
	"github.com/l3uddz/trackarr/config"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
)

/* Privates */
func (r *Release) shouldDelay(pvr *config.PvrInstance) (*int64, error) {
	// iterate delay expressions
	for _, expression := range pvr.DelaysExpr {
		result, err := expr.Run(expression, r)
		if err != nil {
			return nil, errors.Wrapf(err, "failed checking delay expression")
		}

		expResult, ok := result.(int64)
		if !ok {
			return nil, errors.New("failed type asserting delay expression result")
		}

		if expResult > 0 {
			r.Log.Tracef("Delaying release for pvr %q due to delay expression match", pvr.Config.Name)
			return &expResult, nil
		}
	}

	return nil, nil
}
