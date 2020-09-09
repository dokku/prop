package command

import (
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

const (
	// EnvPropCLINoColor is an env var that toggles colored UI output.
	EnvPropCLINoColor = `PROP_CLI_NO_COLOR`
)

// NamedCommand is a interface to denote a commmand's name.
type NamedCommand interface {
	Name() string
}

// Commands returns the mapping of CLI commands for prop. The meta
// parameter lets you set meta options for all commands.
func Commands(metaPtr *Meta, agentUi cli.Ui) map[string]cli.CommandFactory {
	if metaPtr == nil {
		metaPtr = new(Meta)
	}

	meta := *metaPtr
	if meta.Ui == nil {
		meta.Ui = &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      colorable.NewColorableStdout(),
			ErrorWriter: colorable.NewColorableStderr(),
		}
	}

	all := map[string]cli.CommandFactory{}

	for k, v := range KeyValueCommands(meta) {
		all[k] = v
	}

	for k, v := range ListCommands(meta) {
		all[k] = v
	}

	return all
}

func KeyValueCommands(meta Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"del": func() (cli.Command, error) {
			return &DelCommand{Meta: meta}, nil
		},
		"exists": func() (cli.Command, error) {
			return &ExistsCommand{Meta: meta}, nil
		},
		"get": func() (cli.Command, error) {
			return &GetCommand{Meta: meta}, nil
		},
		"get-all": func() (cli.Command, error) {
			return &GetAllCommand{Meta: meta}, nil
		},
		"set": func() (cli.Command, error) {
			return &SetCommand{Meta: meta}, nil
		},
	}
}

func ListCommands(meta Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"lindex": func() (cli.Command, error) {
			return &LindexCommand{Meta: meta}, nil
		},
		"lismember": func() (cli.Command, error) {
			return &LismemberCommand{Meta: meta}, nil
		},
		"llen": func() (cli.Command, error) {
			return &LlenCommand{Meta: meta}, nil
		},
		"lrange": func() (cli.Command, error) {
			return &LrangeCommand{Meta: meta}, nil
		},
		"lrem": func() (cli.Command, error) {
			return &LremCommand{Meta: meta}, nil
		},
		"lset": func() (cli.Command, error) {
			return &LsetCommand{Meta: meta}, nil
		},
		"rpush": func() (cli.Command, error) {
			return &RpushCommand{Meta: meta}, nil
		},
	}
}
