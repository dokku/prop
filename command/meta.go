package command

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/colorstring"
	"github.com/posener/complete"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	// Constants for CLI identifier length
	shortId = 8
	fullId  = 36
)

// FlagSetFlags is an enum to define what flags are present in the
// default FlagSet returned by Meta.FlagSet.
type FlagSetFlags uint

const (
	FlagSetNone    FlagSetFlags = 0
	FlagSetClient  FlagSetFlags = 1 << iota
	FlagSetDefault              = FlagSetClient
)

// Meta contains the meta-options and functionality that nearly every
// Prop command inherits.
type Meta struct {
	Ui cli.Ui

	// Whether to not-colorize output
	noColor bool

	// Namespace to read from
	namespace string

	// URL to read/write data to
	url string
}

func (m *Meta) Namespace() string {
	return m.namespace
}

func (m *Meta) URL() string {
	return m.url
}

// FlagSet returns a FlagSet with the common flags that every
// command implements. The exact behavior of FlagSet can be configured
// using the flags as the second parameter, for example to disable
// server settings on the commands that don't talk to a server.
func (m *Meta) FlagSet(n string, fs FlagSetFlags) *flag.FlagSet {
	f := flag.NewFlagSet(n, flag.ContinueOnError)

	// FlagSetClient is used to enable the settings for specifying
	// client connectivity options.
	if fs&FlagSetClient != 0 {
		f.BoolVar(&m.noColor, "no-color", false, "")
		f.StringVar(&m.namespace, "namespace", "default", "Namespace to use")
		f.StringVar(&m.url, "url", "file:/var/lib/prop/data?system-user=root&system-group=root", "URL to read/write data to")
	}

	f.SetOutput(&uiErrorWriter{ui: m.Ui})

	return f
}

// AutocompleteFlags returns a set of flag completions for the given flag set.
func (m *Meta) AutocompleteFlags(fs FlagSetFlags) complete.Flags {
	if fs&FlagSetClient == 0 {
		return nil
	}

	return complete.Flags{
		"-no-color":  complete.PredictNothing,
		"-namespace": complete.PredictNothing,
		"-url":       complete.PredictNothing,
	}
}

func (m *Meta) Colorize() *colorstring.Colorize {
	return &colorstring.Colorize{
		Colors:  colorstring.DefaultColors,
		Disable: m.noColor || !terminal.IsTerminal(int(os.Stdout.Fd())),
		Reset:   true,
	}
}

// generalOptionsUsage returns the help string for the global options.
func generalOptionsUsage() string {
	helpText := `
  --no-color
    Disables colored command output. Alternatively, PROP_CLI_NO_COLOR may be
    set.
  --namespace <namespace>
    The namespace to interact with.
  --url <url>
    The url to use for the backend.
`
	return strings.TrimSpace(helpText)
}

// funcVar is a type of flag that accepts a function that is the string given
// by the user.
type funcVar func(s string) error

func (f funcVar) Set(s string) error { return f(s) }
func (f funcVar) String() string     { return "" }
func (f funcVar) IsBoolFlag() bool   { return false }

func exampleString(examples map[string]string) string {
	exampleString := []string{}

	for name, example := range examples {
		exampleString = append(exampleString, "  "+name+"\n    $ "+example)
	}

	return strings.Join(exampleString, "\n\n")
}

func flagString(flags *flag.FlagSet) string {
	flagString := []string{}

	flags.VisitAll(func(f *flag.Flag) {
		if f.DefValue == "true" || f.DefValue == "false" {
			flagString = append(flagString, fmt.Sprintf("--%s", f.Name))
			return
		}

		flagString = append(flagString, fmt.Sprintf("--%s <%[1]s-value>", f.Name))
	})
	return strings.Join(flagString, " ")
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
