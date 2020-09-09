package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type SmembersCommand struct {
	Meta
}

func (c *SmembersCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *SmembersCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *SmembersCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *SmembersCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *SmembersCommand) Examples() map[string]string {
	return map[string]string{
		"Get all the members in a set": "prop smembers myset",
	}
}

func (c *SmembersCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *SmembersCommand) Name() string {
	return "smembers"
}

func (c *SmembersCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *SmembersCommand) Synopsis() string {
	return "Get all the members in a set"
}

func (c *SmembersCommand) Run(args []string) int {
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
	members, err := b.Smembers(key)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	for member := range members {
		c.Ui.Output(member)
	}

	return 0
}
