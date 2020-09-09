package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type LindexCommand struct {
	Meta
}

func (c *LindexCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *LindexCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "index",
		Optional: false,
		Type:     ArgumentInt,
	})
	return args
}

func (c *LindexCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *LindexCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LindexCommand) Examples() map[string]string {
	return map[string]string{
		"Get the first element in a list": "prop lindex mykey 0",
	}
}

func (c *LindexCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *LindexCommand) Name() string {
	return "lindex"
}

func (c *LindexCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *LindexCommand) Synopsis() string {
	return "Get an element from a list by its index"
}

func (c *LindexCommand) Run(args []string) int {
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
	index := arguments["index"].IntValue()
	value, err := b.Lindex(key, index)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(value)
	return 0
}
