package config

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PvrConfig struct {
	Name    string
	Enabled bool
	URL     string
	ApiKey  string
	Timeout int
	Filters *PvrFilters
}

type PvrFilters struct {
	Ignores []string
	Accepts []string
	Delays  []string
}

type PvrInstance struct {
	Config             *PvrConfig
	IgnoresExpr        []*vm.Program
	AcceptsExpr        []*vm.Program
	DelaysExpr         []*vm.Program
	HasFileExpressions bool
}

/* Public */

func (p *PvrInstance) ShouldIgnore(r *ReleaseInfo, log *logrus.Entry) (bool, error) {
	// iterate ignore expressions
	for _, expression := range p.IgnoresExpr {
		result, err := expr.Run(expression, r)
		if err != nil {
			return true, errors.Wrap(err, "failed checking ignore expression")
		}

		expResult, ok := result.(bool)
		if !ok {
			return true, errors.New("failed type asserting ignore expression result")
		}

		if expResult {
			if log != nil {
				log.Tracef("Ignoring release for pvr %q due to ignore expression match", p.Config.Name)
			}

			return true, nil
		}
	}

	return false, nil
}

func (p *PvrInstance) ShouldAccept(r *ReleaseInfo, log *logrus.Entry) (bool, error) {
	// iterate accept expressions
	for _, expression := range p.AcceptsExpr {
		result, err := expr.Run(expression, r)
		if err != nil {
			return false, errors.Wrap(err, "failed checking accept expression")
		}

		expResult, ok := result.(bool)
		if !ok {
			return false, errors.New("failed type asserting accept expression result")
		}

		if expResult {
			if log != nil {
				log.Tracef("Allowing release for pvr %q due to accept expression match", p.Config.Name)
			}

			return true, nil
		}
	}

	return false, nil
}

func (p *PvrInstance) ShouldDelay(r *ReleaseInfo, log *logrus.Entry) (*int64, error) {
	// iterate delay expressions
	for _, expression := range p.DelaysExpr {
		result, err := expr.Run(expression, r)
		if err != nil {
			return nil, errors.Wrap(err, "failed checking delay expression")
		}

		expResult, ok := result.(int64)
		if !ok {
			return nil, errors.New("failed type asserting delay expression result")
		}

		if expResult > 0 {
			if log != nil {
				log.Tracef("Delaying release for pvr %q due to delay expression match", p.Config.Name)
			}

			return &expResult, nil
		}
	}

	return nil, nil
}
