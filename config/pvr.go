package config

import (
	"github.com/antonmedv/expr/vm"
)

type PvrConfig struct {
	Name    string
	Enabled bool
	URL     string
	ApiKey  string
	Filters *PvrFilters
}

type PvrFilters struct {
	Ignores []string
	Accepts []string
	Delays  []string
}

type PvrInstance struct {
	Config      *PvrConfig
	IgnoresExpr []*vm.Program
	AcceptsExpr []*vm.Program
	DelaysExpr  []*vm.Program
}
