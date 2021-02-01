package cli

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

const (
	optionStringModuleDescribe = "describe"
	optionStringModuleStart    = "start"
	optionStringModuleStop     = "stop"
	optionStringModuleEnable   = "enable"
	optionStringModuleDisable  = "disable"
)

func (cli CLI) getModuleSuggestions() (s []prompt.Suggest) {
	for _, m := range cli.modules {
		s = append(s, prompt.Suggest{Text: m.Name, Description: ""})
	}
	return
}

func (cli CLI) describeModule(s string) {
	m := cli.modules[s]
	if m != nil {
		fmt.Printf("%v\n", m)
	}
}

func (cli CLI) startModule(s string) {
	m := cli.modules[s]
	if m != nil {
		m.Start()
	}
}

func (cli CLI) stopModule(s string) {
	m := cli.modules[s]
	if m != nil {
		m.Stop()
	}
}

func (cli CLI) enableModule(s string) {
	m := cli.modules[s]
	if m != nil {
		m.Enabled = true
	}
}

func (cli CLI) disableModule(s string) {
	m := cli.modules[s]
	if m != nil {
		m.Enabled = false
	}
}

func (cli CLI) modulesRootHandler(s string) {
	switch s {
	case optionStringModuleDescribe:
		cli.modules.Describe()
	case optionStringModuleStart:
		cli.modules.Start()
	case optionStringModuleStop:
		cli.modules.Stop()
	case optionStringModuleEnable:
		cli.modules.Enable()
	case optionStringModuleDisable:
		cli.modules.Disable()
	default:
		fmt.Printf("Unknown modules command: %s\n", s)
	}
}

func (cli CLI) newModuleOption() (optionModule option) {
	optionModule = option{
		func() []prompt.Suggest {
			return []prompt.Suggest{
				{Text: optionStringModuleDescribe, Description: ""},
				{Text: optionStringModuleStart, Description: ""},
				{Text: optionStringModuleStop, Description: ""},
				{Text: optionStringModuleEnable, Description: ""},
				{Text: optionStringModuleDisable, Description: ""},
			}
		}, options{
			optionStringModuleDescribe: option{cli.getModuleSuggestions, options{}, cli.describeModule},
			optionStringModuleStart:    option{cli.getModuleSuggestions, options{}, cli.startModule},
			optionStringModuleStop:     option{cli.getModuleSuggestions, options{}, cli.stopModule},
			optionStringModuleEnable:   option{cli.getModuleSuggestions, options{}, cli.enableModule},
			optionStringModuleDisable:  option{cli.getModuleSuggestions, options{}, cli.disableModule},
		}, cli.modulesRootHandler,
	}
	return
}
