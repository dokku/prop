package command

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type GetAllCommand struct {
	Meta
}

func (c *GetAllCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *GetAllCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "prefix",
		Optional: true,
		Type:     ArgumentString,
	})
	return args
}

func (c *GetAllCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *GetAllCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *GetAllCommand) Examples() map[string]string {
	return map[string]string{
		"Get all values in a namespace": "prop get-all",
	}
}

func (c *GetAllCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *GetAllCommand) Name() string {
	return "get-all"
}

func (c *GetAllCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *GetAllCommand) Synopsis() string {
	return "Get all values in a namespace"
}

func (c *GetAllCommand) Run(args []string) int {
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

	var keyValuePairs map[string]string
	if arguments["prefix"].HasValue {
		keyValuePairs, err = b.GetAllByPrefix(arguments["prefix"].StringValue())
	} else {
		keyValuePairs, err = b.GetAll()
	}

	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	var kv []string
	for key, value := range keyValuePairs {
		kv = append(kv, fmt.Sprintf("%v | %v", key, value))
	}

	if len(kv) == 0 {
		return 0
	}

	c.Ui.Output(formatKV(kv))

	return 0
}
