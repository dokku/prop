package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type NamespaceClearCommand struct {
	Meta
}

func (c *NamespaceClearCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *NamespaceClearCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "namespace",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *NamespaceClearCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *NamespaceClearCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *NamespaceClearCommand) Examples() map[string]string {
	return map[string]string{
		"Delete keys in a namespace": "prop namespace clear mynamespace",
	}
}

func (c *NamespaceClearCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *NamespaceClearCommand) Name() string {
	return "namespace clear"
}

func (c *NamespaceClearCommand) Synopsis() string {
	return "Delete all keys from a given namespace"
}

func (c *NamespaceClearCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *NamespaceClearCommand) Run(args []string) int {
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

	namespace := arguments["namespace"].StringValue()
	success, err := b.NamespaceClear(namespace)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !success {
		return 1
	}

	return 1
}
