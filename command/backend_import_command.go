package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type BackendImportCommand struct {
	Meta

	clearBackend bool
}

func (c *BackendImportCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *BackendImportCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "path",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *BackendImportCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *BackendImportCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *BackendImportCommand) Examples() map[string]string {
	return map[string]string{
		"Import a property collection": "prop backend import /tmp/backend.json",
	}
}

func (c *BackendImportCommand) FlagSet() *flag.FlagSet {
	f := c.Meta.FlagSet(c.Name(), FlagSetClient)
	f.BoolVar(&c.clearBackend, "clear-backend", false, "")
	return f
}

func (c *BackendImportCommand) Name() string {
	return "backend import"
}

func (c *BackendImportCommand) Synopsis() string {
	return `Import a backend to a json file

  When importing a backend, properties are merged into the existing backend
  unless the --clear flag is specified.

  When migrating a backend, it is assumed that there are is no concurrent
  access to the backend. In other words, if another process is changing
  values of the backend, then the import may result in an invalid state.`
}

func (c *BackendImportCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *BackendImportCommand) Run(args []string) int {
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

	path := arguments["path"].StringValue()
	p, err := backend.DeserializePropertyCollection(path)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	imported, err := b.BackendImport(p, c.clearBackend)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !imported {
		return 1
	}

	return 1
}
