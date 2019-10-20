package release

import (
	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/l3uddz/trackarr/autodl/parser"
	"github.com/l3uddz/trackarr/config"
	"github.com/l3uddz/trackarr/logger"
	"github.com/pkg/errors"
)

var (
	log            = logger.GetLogger("release")
	pvrExpressions map[*config.PvrConfiguration]map[string][]*vm.Program
)

/* Public */

func Init(pvr *[]config.PvrConfiguration, cfg *config.TrackerConfiguration) error {
	// placeholder release for expression compilation
	tr := &TrackerRelease{
		Tracker: &parser.TrackerInfo{
			LongName:          "Broadcasthenet",
			ShortName:         nil,
			Settings:          nil,
			Servers:           nil,
			Channels:          nil,
			Announcers:        nil,
			IgnoreLines:       nil,
			LinePatterns:      nil,
			MultiLinePatterns: nil,
			LineMatchedRules:  nil,
		},
		Log:         log,
		Cfg:         cfg,
		TrackerName: "BTN",
		ReleaseTime: "2014-02-10T00:00:00Z",
		TorrentName: "Batwoman.S01E01.Pilot.1080p.AMZN.WEB-DL.DDP5.1.H.264-NTb",
		TorrentURL:  "https://www.google.com",
		SizeString:  "2.63 GB",
		SizeBytes:   2834678416,
		Category:    "Episode",
		Encoder:     "x264",
		Resolution:  "1080p",
		Container:   "MKV",
		Origin:      "Scene",
		Tags:        "BD25 , Blu-ray , m2ts , 1080p , Scene",
	}

	// iterate pvrs
	enabledPvrs := 0
	pvrExpressions = make(map[*config.PvrConfiguration]map[string][]*vm.Program, 0)
	for _, obj := range *pvr {
		pvrObj := obj

		// skip disabled pvr
		if !pvrObj.Enabled {
			continue
		} else {
			enabledPvrs++
		}

		ignoreExpressions := make([]*vm.Program, 0)
		acceptExpressions := make([]*vm.Program, 0)

		// iterate pvr ignore expressions
		compiledIgnoreExpressions := 0
		for _, ignoreExpression := range pvrObj.Ignores {
			program, err := expr.Compile(ignoreExpression, expr.Env(tr))
			if err != nil {
				// failed to compile expression, return error
				return errors.Wrapf(err, "failed compiling ignore expression for pvr: %q", pvrObj.Name)
			}

			ignoreExpressions = append(ignoreExpressions, program)
			compiledIgnoreExpressions++
		}

		// iterate pvr accept expressions
		compiledAcceptExpressions := 0
		for _, acceptExpression := range pvrObj.Accepts {
			program, err := expr.Compile(acceptExpression, expr.Env(tr))
			if err != nil {
				// failed to compile expression, return error
				return errors.Wrapf(err, "failed compiling accept expression for pvr: %q", pvrObj.Name)
			}

			acceptExpressions = append(acceptExpressions, program)
			compiledAcceptExpressions++
		}

		// store compiled expressions (even if nothing was set - storing the pvr object is required for iterating later)
		pvrExpressions[&pvrObj] = map[string][]*vm.Program{
			"ignores": ignoreExpressions,
			"accepts": acceptExpressions,
		}

		log.Debugf("Compiled expressions for pvr %q: %d ignores, %d accepts", pvrObj.Name, compiledIgnoreExpressions,
			compiledAcceptExpressions)
	}

	if enabledPvrs == 0 {
		return errors.New("no pvr's were enabled")
	}

	return nil
}
