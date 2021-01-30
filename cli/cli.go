package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/jeffssh/Ophanim/module"
)

type getSuggestionsFunc func() []prompt.Suggest

type option struct {
	getSuggestions getSuggestionsFunc
	opts           options
	//do          executor
}

type options map[string]option

// CLI - cli interface for Ophanim
type CLI struct {
	modules module.Map
	prompt  *prompt.Prompt
	opt     option
}

var (
	promptText        = "( o) "
	optExit           = "exit"
	optModule         = "modules"
	optModuleDescribe = "describe"
	optModuleStart    = "start"
	optModuleStop     = "stop"
)

// New - initialize the CLI
func New() (cli CLI) {
	cli.modules = module.LoadAllModules("./module/modules/yaml/")
	var moduleOpt = option{
		func() []prompt.Suggest {
			return []prompt.Suggest{
				{Text: optModuleDescribe, Description: ""},
				{Text: optModuleStart, Description: ""},
				{Text: optModuleStop, Description: ""},
			}
		}, options{
			optModuleDescribe: option{cli.getModuleSuggestions(), options{}},
			optModuleStart:    option{cli.getModuleSuggestions(), options{}},
			optModuleStop:     option{cli.getModuleSuggestions(), options{}},
		},
	}

	var exitOpt = option{
		func() []prompt.Suggest {
			return []prompt.Suggest{
				{Text: optExit, Description: ""},
			}
		}, options{},
	}

	var opt = option{
		func() []prompt.Suggest {
			return []prompt.Suggest{}
		},
		options{
			"": option{
				func() []prompt.Suggest {
					return []prompt.Suggest{
						{Text: optModule, Description: ""},
						{Text: optExit, Description: ""},
					}
				}, options{
					optModule: moduleOpt,
					optExit:   exitOpt,
				},
			},
		},
	}
	cli.opt = opt
	cli.prompt = prompt.New(
		cli.test,
		cli.completer,
		prompt.OptionPrefix(promptText),
		prompt.OptionPrefixTextColor(prompt.Red),
		prompt.OptionSuggestionBGColor(prompt.Red),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkRed),
	)
	setupCloseHandler(cli.modules)
	return
}

func (cli CLI) completer(d prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{}
	if d.TextBeforeCursor() == "" {
		return s
	}

	args := strings.Split(d.TextBeforeCursor(), " ")
	args = append([]string{""}, args[:len(args)-1]...)
	o := cli.opt
	for _, arg := range args {
		o = o.opts[arg]
		if &o == nil {
			return s
		}
	}
	if o.getSuggestions != nil {
		s = o.getSuggestions()
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

// Prompt - display the main prompt
func (cli CLI) Prompt() {
	cli.prompt.Run()
}

func (cli CLI) getModuleSuggestions() getSuggestionsFunc {
	return func() (s []prompt.Suggest) {
		for _, m := range cli.modules {
			s = append(s, prompt.Suggest{Text: m.Name, Description: ""})
		}
		return
	}
}

func (cli CLI) test(s string) {
	opt := s
	switch opt {
	case optModule:
		cli.modules.Start()
	case optExit:
		cli.modules.Stop()
		os.Exit(0)
	default:
		fmt.Printf("Unknown syntax: %s\n", opt)
	}
}

func setupCloseHandler(modules module.Map) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("\rCtrl+C pressed in terminal, stopping modules")
		modules.Stop()
		os.Exit(0)
	}()
}
