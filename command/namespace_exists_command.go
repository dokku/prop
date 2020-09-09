package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type NamespaceExistsCommand struct {
	Meta
}

func (c *NamespaceExistsCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *NamespaceExistsCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "namespace",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *NamespaceExistsCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *NamespaceExistsCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *NamespaceExistsCommand) Examples() map[string]string {
	return map[string]string{
		"Check if a namespace exists": "prop namespace exists mynamespace",
	}
}

func (c *NamespaceExistsCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *NamespaceExistsCommand) Name() string {
	return "namespace exists"
}

func (c *NamespaceExistsCommand) Synopsis() string {
	return "Checks if there are any keys in a given namespace"
}

func (c *NamespaceExistsCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *NamespaceExistsCommand) Run(args []string) int {
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
	exists, err := b.NamespaceExists(namespace)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !exists {
		return 1
	}

	return 1
}
