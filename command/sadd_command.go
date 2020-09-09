package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type SaddCommand struct {
	Meta
}

func (c *SaddCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *SaddCommand) Arguments() []Argument {
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

func (c *SaddCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *SaddCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *SaddCommand) Examples() map[string]string {
	return map[string]string{
		"Add a member to the set": "prop sadd myset mymember",
	}
}

func (c *SaddCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *SaddCommand) Name() string {
	return "sadd"
}

func (c *SaddCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *SaddCommand) Synopsis() string {
	return "Add one or more members to a set"
}

func (c *SaddCommand) Run(args []string) int {
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
	addedCount, err := b.Sadd(key, members...)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output(fmt.Sprintf("%d", addedCount))

	return 0
}
