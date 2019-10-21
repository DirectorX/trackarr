package pvr

import (
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
)

var (
	log = logger.GetLogger("pvr")
)

/* Public */

func Init() error {
	for _, p := range config.Config.Pvr {
		// skip disabled trackers
		if !p.Enabled {
			log.Debugf("Skipping disabled PVR: %s", p.Name)

			continue
		}

		config.Pvr[p.Name] = &config.PvrInstance{
			Config: &p,
		}

		// Compile expressions
		if err := compileExpr(config.Pvr[p.Name]); err != nil {
			return err
		}

		log.Infof("Initialized PVR %s", p.Name)
	}

	return nil
}

func compileExpr(p *config.PvrInstance) error {
	exprEnv := &config.ReleaseInfo{}

	// iterate pvr ignore expressions
	for _, ignoreExpr := range p.Config.Ignores {
		program, err := expr.Compile(ignoreExpr, expr.Env(exprEnv), expr.AsBool())
		if err != nil {
			return errors.Wrapf(err, "failed compiling ignore expression for pvr: %q", p.Config.Name)
		}

		p.IgnoresExpr = append(p.IgnoresExpr, program)
	}

	// iterate pvr accept expressions
	for _, acceptExpr := range p.Config.Accepts {
		program, err := expr.Compile(acceptExpr, expr.Env(exprEnv), expr.AsBool())
		if err != nil {
			return errors.Wrapf(err, "failed compiling accept expression for pvr: %q", p.Config.Name)
		}

		p.AcceptsExpr = append(p.AcceptsExpr, program)
	}

	// iterate pvr delay expressions
	for _, delayExpr := range p.Config.Delays {
		program, err := expr.Compile(delayExpr, expr.Env(exprEnv), expr.AsInt64())
		if err != nil {
			return errors.Wrapf(err, "failed compiling delay expression for pvr: %q", p.Config.Name)
		}

		p.DelaysExpr = append(p.DelaysExpr, program)
	}

	log.Debugf("Compiled expressions for pvr %q: %d ignores, %d accepts, %d delays",
		p.Config.Name,
		len(p.IgnoresExpr),
		len(p.AcceptsExpr),
		len(p.DelaysExpr),
	)

	return nil
}
