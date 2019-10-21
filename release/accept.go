package release

import (
	"github.com/l3uddz/trackarr/config"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/pkg/errors"
)

/* Privates */
func (r *Release) shouldAccept(pvr *config.PvrConfig, expressions *map[string][]*vm.Program) (bool, error) {
	acceptExpressions, ok := (*expressions)["accepts"]
	if !ok {
		// there were no accepts
		return true, nil
	}

	// iterate accept expressions
	for _, expression := range acceptExpressions {
		result, err := expr.Run(expression, r)
		if err != nil {
			return false, errors.Wrapf(err, "failed checking accept expression")
		}

		expResult, ok := result.(bool)
		if !ok {
			return false, errors.New("failed type asserting accept expression result")
		}

		if expResult {
			r.Log.Tracef("Allowing release for pvr %q due to accept expression match", pvr.Name)
			return true, nil
		}
	}

	return false, nil
}
