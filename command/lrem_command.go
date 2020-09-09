package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type LremCommand struct {
	Meta
}

func (c *LremCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *LremCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "count",
		Optional: false,
		Type:     ArgumentInt,
	})
	args = append(args, Argument{
		Name:     "element",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *LremCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *LremCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LremCommand) Examples() map[string]string {
	return map[string]string{
		"Remove an element from a list": "prop lrem mykey 0 myelement",
	}
}

func (c *LremCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *LremCommand) Name() string {
	return "lrem"
}

func (c *LremCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *LremCommand) Synopsis() string {
	return "Remove elements from a list"
}

func (c *LremCommand) Run(args []string) int {
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
	count := arguments["count"].IntValue()
	element := arguments["element"].StringValue()
	removedCount, err := b.Lrem(key, count, element)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("%d", removedCount))

	return 0
}
