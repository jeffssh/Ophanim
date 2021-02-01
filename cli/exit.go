package cli

import (
	"fmt"
	"os"

	"github.com/c-bata/go-prompt"
)

const (
	optionStringExit = "exit"
)

func (cli CLI) exitRootHandler(s string) {
	fmt.Printf("\rExiting, stopping modules\n")
	cli.modules.Stop()
	os.Exit(0)
}

func (cli CLI) newExitOption() (optionExit option) {
	optionExit = option{
		func() []prompt.Suggest {
			return []prompt.Suggest{
				{Text: optionStringExit, Description: ""},
			}
		}, options{}, cli.exitRootHandler,
	}
	return
}
