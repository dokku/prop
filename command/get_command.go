package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type GetCommand struct {
	Meta
}

func (c *GetCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *GetCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "default-value",
		Optional: true,
		Type:     ArgumentString,
	})
	return args
}

func (c *GetCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *GetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *GetCommand) Examples() map[string]string {
	return map[string]string{
		"Get a key": "prop get mykey",
	}
}
func (c *GetCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *GetCommand) Name() string {
	return "get"
}

func (c *GetCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *GetCommand) Synopsis() string {
	return "Get the value of a key"
}

func (c *GetCommand) Run(args []string) int {
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
	defaultValue := arguments["default-value"].StringValue()
	value, err := b.Get(key, defaultValue)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(value)
	return 0
}
