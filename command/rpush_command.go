package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type RpushCommand struct {
	Meta
}

func (c *RpushCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *RpushCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "elements",
		Optional: false,
		Type:     ArgumentList,
	})
	return args
}

func (c *RpushCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *RpushCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RpushCommand) Examples() map[string]string {
	return map[string]string{
		"Add elements to the end of a list": "prop rpush mykey myelement mysecondelement",
	}
}

func (c *RpushCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *RpushCommand) Name() string {
	return "rpush"
}

func (c *RpushCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *RpushCommand) Synopsis() string {
	return "Append one or more elements to a list"
}

func (c *RpushCommand) Run(args []string) int {
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
	elements := arguments["elements"].ListValue()
	length, err := b.Rpush(key, elements...)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("%d", length))

	return 0
}
