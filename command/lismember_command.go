package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type LismemberCommand struct {
	Meta
}

func (c *LismemberCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *LismemberCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "element",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *LismemberCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *LismemberCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LismemberCommand) Examples() map[string]string {
	return map[string]string{
		"Check if an element is in the list": "prop lismember mykey myelement",
	}
}

func (c *LismemberCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *LismemberCommand) Name() string {
	return "lismember"
}

func (c *LismemberCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *LismemberCommand) Synopsis() string {
	return "Determine if a given value is an element in the list"
}

func (c *LismemberCommand) Run(args []string) int {
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
	element := arguments["element"].StringValue()
	ok, err := b.Lismember(key, element)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !ok {
		return 1
	}

	return 0
}
