package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type LsetCommand struct {
	Meta
}

func (c *LsetCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *LsetCommand) Arguments() []Argument {
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
	args = append(args, Argument{
		Name:     "element",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *LsetCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *LsetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LsetCommand) Examples() map[string]string {
	return map[string]string{
		"Set an element at a position in a list": "prop lset mylist 0 myelement",
	}
}

func (c *LsetCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *LsetCommand) Name() string {
	return "lset"
}

func (c *LsetCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *LsetCommand) Synopsis() string {
	return "Set the value of an element in a list by its index"
}

func (c *LsetCommand) Run(args []string) int {
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
	element := arguments["element"].StringValue()
	ok, err := b.Lset(key, index, element)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !ok {
		return 1
	}

	return 0
}
