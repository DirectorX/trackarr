package pvr

import (
	"fmt"
	"gitlab.com/cloudb0x/trackarr/config"
	"gitlab.com/cloudb0x/trackarr/logger"
	"strings"

	"github.com/antonmedv/expr"
	"github.com/pkg/errors"
)

var (
	log = logger.GetLogger("pvr")
)

/* Public */

func Init() error {
	for _, p := range config.Config.Pvr {
		// skip disabled pvr
		if !p.Enabled {
			log.Debugf("Skipping disabled PVR: %s", p.Name)
			continue
		}

		// check if pvr has already been loaded (duplicate pvr name)
		if _, exists := config.Pvr[p.Name]; exists {
			return fmt.Errorf("pvr with the same name already loaded: %q", p.Name)
		}

		// init pvr instance
		p2 := p
		config.Pvr[p.Name] = &config.PvrInstance{
			Config: &p2,
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

	// dont compile filters when not set
	if p.Config.Filters != nil {

		// iterate pvr ignore expressions
		for _, ignoreExpr := range p.Config.Filters.Ignores {
			program, err := expr.Compile(ignoreExpr, expr.Env(exprEnv), expr.AsBool())
			if err != nil {
				return errors.Wrapf(err, "failed compiling ignore expression for pvr: %q", p.Config.Name)
			}

			p.IgnoresExpr = append(p.IgnoresExpr, program)

			if strings.Contains(ignoreExpr, "Files") {
				// used by proxy torrent endpoint to check expressions (only ignore will trigger second-sweep)
				p.HasFileExpressions = true
			}
		}

		// iterate pvr accept expressions
		for _, acceptExpr := range p.Config.Filters.Accepts {
			program, err := expr.Compile(acceptExpr, expr.Env(exprEnv), expr.AsBool())
			if err != nil {
				return errors.Wrapf(err, "failed compiling accept expression for pvr: %q", p.Config.Name)
			}

			p.AcceptsExpr = append(p.AcceptsExpr, program)
		}

		// iterate pvr delay expressions
		for _, delayExpr := range p.Config.Filters.Delays {
			program, err := expr.Compile(delayExpr, expr.Env(exprEnv), expr.AsInt64())
			if err != nil {
				return errors.Wrapf(err, "failed compiling delay expression for pvr: %q", p.Config.Name)
			}

			p.DelaysExpr = append(p.DelaysExpr, program)
		}
	}

	log.Debugf("Compiled expressions for pvr %q: %d ignores, %d accepts, %d delays",
		p.Config.Name,
		len(p.IgnoresExpr),
		len(p.AcceptsExpr),
		len(p.DelaysExpr),
	)

	return nil
}
