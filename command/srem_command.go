package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type SremCommand struct {
	Meta
}

func (c *SremCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *SremCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "members",
		Optional: false,
		Type:     ArgumentList,
	})
	return args
}

func (c *SremCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *SremCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *SremCommand) Examples() map[string]string {
	return map[string]string{
		"Remove a members from a set": "prop srem myset mymember",
	}
}

func (c *SremCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *SremCommand) Name() string {
	return "srem"
}

func (c *SremCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *SremCommand) Synopsis() string {
	return "Remove one or more members from a set"
}

func (c *SremCommand) Run(args []string) int {
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
	members := arguments["members"].ListValue()
	removedCount, err := b.Srem(key, members...)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("%d", removedCount))

	return 0
}
