package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type BackendResetCommand struct {
	Meta

	clearBackend bool
}

func (c *BackendResetCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *BackendResetCommand) Arguments() []Argument {
	args := []Argument{}
	return args
}

func (c *BackendResetCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *BackendResetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *BackendResetCommand) Examples() map[string]string {
	return map[string]string{
		"Reset a backend": "prop backend reset",
	}
}

func (c *BackendResetCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *BackendResetCommand) Name() string {
	return "backend reset"
}

func (c *BackendResetCommand) Synopsis() string {
	return "Clear all values in a backend"
}

func (c *BackendResetCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *BackendResetCommand) Run(args []string) int {
	flags := c.FlagSet()
	flags.Usage = func() { c.Ui.Output(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	_, err := c.ParsedArguments(flags.Args())
	if err != nil {
		c.Ui.Error(err.Error())
		c.Ui.Error(commandErrorText(c))
		return 1
	}

	b, err := backend.ConstructBackend(c.Meta.URL(), c.Meta.Namespace())
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	success, err := b.BackendReset()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !success {
		return 1
	}

	return 1
}
