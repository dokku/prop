package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type BackendExportCommand struct {
	Meta
}

func (c *BackendExportCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *BackendExportCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "path",
		Optional: false,
		Type:     ArgumentString,
	})
	return args
}

func (c *BackendExportCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *BackendExportCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *BackendExportCommand) Examples() map[string]string {
	return map[string]string{
		"Export a property collection": "prop backend export /tmp/backend.json",
	}
}

func (c *BackendExportCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *BackendExportCommand) Name() string {
	return "backend export"
}

func (c *BackendExportCommand) Synopsis() string {
	return `Exports a backend to a json file

  When export a backend, it is assumed that there are is no concurrent access
  to the backend. In other words, if another process is changing values of the
  backend, then the export may result in an invalid state.`
}

func (c *BackendExportCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *BackendExportCommand) Run(args []string) int {
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

	p, err := b.BackendExport()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	path := arguments["path"].StringValue()
	success, err := backend.SerializePropertyCollection(p, path)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if !success {
		return 1
	}

	return 1
}
