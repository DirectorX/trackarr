package release

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/l3uddz/trackarr/config"
	"github.com/pkg/errors"
)

/* Privates */
func (r *TrackerRelease) shouldIgnore(pvr *config.PvrConfiguration, expressions *map[string][]*vm.Program) (bool, error) {
	ignoreExpressions, ok := (*expressions)["ignores"]
	if !ok {
		// there were no ignores
		return false, nil
	}

	// iterate ignore expressions
	for _, expression := range ignoreExpressions {
		result, err := expr.Run(expression, r)
		if err != nil {
			return true, errors.Wrapf(err, "failed checking ignore expression")
		}

		expResult, ok := result.(bool)
		if !ok {
			return true, errors.New("failed type asserting ignore expression result")
		}

		if expResult {
			r.Log.Tracef("Ignoring release for pvr %q due to ignore expression match", pvr.Name)
			return true, nil
		}
	}

	return false, nil
}
