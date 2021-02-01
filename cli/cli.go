package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/c-bata/go-prompt"
	"github.com/jeffssh/Ophanim/module"
)

type option struct {
	getSuggestions func() []prompt.Suggest
	opts           options
	do             func(s string)
}

type options map[string]option

// CLI - cli interface for Ophanim
type CLI struct {
	modules module.Map
	prompt  *prompt.Prompt
	opt     option
}

const (
	promptText         = "( o) "
	optionStringModule = "modules"
)

func nop(s string) {
}

// New - initialize the CLI
func New(moduleYamlDir string) (cli CLI) {
	cli.modules = module.LoadAllModules(moduleYamlDir)

	var optionModule = cli.newModuleOption()
	var optionExit = cli.newExitOption()

	var optionRoot = option{
		func() []prompt.Suggest {
			return []prompt.Suggest{}
		},
		options{
			"": option{
				func() []prompt.Suggest {
					return []prompt.Suggest{
						{Text: optionStringModule, Description: ""},
						{Text: optionStringExit, Description: ""},
					}
				}, options{
					optionStringModule: optionModule,
					optionStringExit:   optionExit,
				}, nop,
			},
		}, nop,
	}
	cli.opt = optionRoot
	cli.prompt = prompt.New(
		cli.executor,
		cli.completer,
		prompt.OptionPrefix(promptText),
		prompt.OptionPrefixTextColor(prompt.Red),
		prompt.OptionSuggestionBGColor(prompt.Red),
		prompt.OptionSelectedSuggestionBGColor(prompt.DarkRed),
	)
	cli.setupCloseHandler()
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

func (cli CLI) executor(s string) {
	args := strings.Split(s, " ")
	w, args := args[len(args)-1], args[:len(args)-1]
	args = append([]string{""}, args...)

	o := cli.opt
	for _, arg := range args {
		o = o.opts[arg]
	}

	if o.do != nil {
		o.do(w)
	} else {
		fmt.Printf("Unknown syntax: %s\n", s)
	}
}

func (cli CLI) setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\rCtrl+C pressed in terminal, stopping modules\n")
		cli.modules.Stop()
		os.Exit(0)
	}()
}
