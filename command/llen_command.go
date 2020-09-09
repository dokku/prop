package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type LlenCommand struct {
	Meta
}

func (c *LlenCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *LlenCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *LlenCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *LlenCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LlenCommand) Examples() map[string]string {
	return map[string]string{
		"Get the length of the list": "prop llen mykey",
	}
}

func (c *LlenCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *LlenCommand) Name() string {
	return "llen"
}

func (c *LlenCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *LlenCommand) Synopsis() string {
	return "Get the length of a list"
}

func (c *LlenCommand) Run(args []string) int {
	flags := c.FlagSet()
	flags.Usage = func() { c.Ui.Output(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	arguments, err := c.ParsedArguments(flags.Args())
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

	key := arguments["key"].StringValue()
	length, err := b.Llen(key)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("%d", length))

	return 0
}
