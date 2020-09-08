package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type ExistsCommand struct {
	Meta
}

func (c *ExistsCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

  Check if a key exists:
      $ prop exists mykey
`
	return strings.TrimSpace(helpText)
}

func (c *ExistsCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *ExistsCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *ExistsCommand) Synopsis() string {
	return "Check if a key exists"
}

func (c *ExistsCommand) Name() string {
	return "exists"
}

func (c *ExistsCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *ExistsCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *ExistsCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *ExistsCommand) Run(args []string) int {
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
	ok, err := b.Exists(key)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !ok {
		return 1
	}

	return 0
}
