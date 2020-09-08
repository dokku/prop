package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type SetCommand struct {
	Meta
}

func (c *SetCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

  Set a key:
      $ prop set mykey myvalue
`
	return strings.TrimSpace(helpText)
}

func (c *SetCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *SetCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *SetCommand) Synopsis() string {
	return "Set the value of a key"
}

func (c *SetCommand) Name() string {
	return "set"
}

func (c *SetCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *SetCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "value",
		Optional: true,
		Type:     ArgumentString,
	})
	return args
}

func (c *SetCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *SetCommand) Run(args []string) int {
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
	value := arguments["value"].StringValue()
	ok, err := b.Set(key, value)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !ok {
		return 1
	}

	return 0
}
