package command

import (
	"flag"
	"strings"

	"github.com/dokku/prop/backend"
	"github.com/posener/complete"
)

type LrangeCommand struct {
	Meta
}

func (c *LrangeCommand) Help() string {
	helpText := `
Usage: prop ` + c.Name() + ` ` + flagString(c.FlagSet()) + ` ` + argumentString(c.Arguments()) + `

  ` + c.Synopsis() + `

General Options:
  ` + generalOptionsUsage() + `

Example:

` + exampleString(c.Examples())

	return strings.TrimSpace(helpText)
}

func (c *LrangeCommand) Arguments() []Argument {
	args := []Argument{}
	args = append(args, Argument{
		Name:     "key",
		Optional: false,
		Type:     ArgumentString,
	})
	args = append(args, Argument{
		Name:     "start",
		Optional: true,
		Type:     ArgumentInt,
	})
	args = append(args, Argument{
		Name:     "stop",
		Optional: true,
		Type:     ArgumentInt,
	})
	return args
}

func (c *LrangeCommand) AutocompleteFlags() complete.Flags {
	return complete.Flags{}
}

func (c *LrangeCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *LrangeCommand) Examples() map[string]string {
	return map[string]string{
		"Get all elements in list": "prop lrange mykey",
	}
}

func (c *LrangeCommand) FlagSet() *flag.FlagSet {
	return c.Meta.FlagSet(c.Name(), FlagSetClient)
}

func (c *LrangeCommand) Name() string {
	return "lrange"
}

func (c *LrangeCommand) ParsedArguments(args []string) (map[string]Argument, error) {
	return parseArguments(args, c.Arguments())
}

func (c *LrangeCommand) Synopsis() string {
	return "Get a range of elements from a list"
}

func (c *LrangeCommand) Run(args []string) int {
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

	var values []string
	key := arguments["key"].StringValue()
	start := arguments["start"].IntValue()
	stop := arguments["stop"].IntValue()
	if !arguments["start"].HasValue {
		values, err = b.Lrange(key)
	} else if !arguments["stop"].HasValue {
		values, err = b.Lrangefrom(key, start)
	} else {
		values, err = b.Lrangefromto(key, start, stop)
	}

	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	for _, value := range values {
		c.Ui.Output(value)
	}

	return 0
}
