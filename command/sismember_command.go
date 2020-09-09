package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type SismemberCommand struct {
	Meta
}

func (c *SismemberCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *SismemberCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "member",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *SismemberCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *SismemberCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *SismemberCommand) Examples() map[string]string {
	return map[string]string{
		"Check if an member is in the set": "prop sismember myset mymember",
	}
}

func (c *SismemberCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *SismemberCommand) Name() string {
	return "sismember"
}

func (c *SismemberCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *SismemberCommand) Synopsis() string {
	return "Determine if a given value is a member of a set"
}

func (c *SismemberCommand) Run(args []string) int {
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
	member := arguments["member"].StringValue()
	ok, err := b.Sismember(key, member)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !ok {
		return 1
	}

	return 0
}
