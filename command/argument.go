package command

import (
	"fmt"
	"strings"
)

type Argument struct {
	Name     string
	Optional bool
	Type     ArgumentType
	Value    interface{}
	HasValue bool
}

// ArgumentType is an enum to define what arguments are present
type ArgumentType uint

const (
	ArgumentString ArgumentType = 0
	ArgumentInt    ArgumentType = 1 << iota
	ArgumentBool   ArgumentType = 2 << iota
)

func (a Argument) BoolValue() bool {
	if a.Type != ArgumentBool {
		panic(fmt.Errorf("Unexpected argument type for %s when calling BoolValue()", a.Name))
	}

	return a.Value.(bool)
}

func (a Argument) IntValue() int {
	if a.Type != ArgumentInt {
		panic(fmt.Errorf("Unexpected argument type for %s when calling IntValue()", a.Name))
	}

	return a.Value.(int)
}

func (a Argument) StringValue() string {
	if a.Type != ArgumentString {
		panic(fmt.Errorf("Unexpected argument type for %s when calling StringValue()", a.Name))
	}

	return a.Value.(string)
}

func argumentString(arguments []Argument) string {
	argumentString := []string{}

	for _, argument := range arguments {
		if argument.Optional {
			argumentString = append(argumentString, fmt.Sprintf("[%s]", argument.Name))
		} else {
			argumentString = append(argumentString, fmt.Sprintf("<%s>", argument.Name))
		}
	}

	return strings.Join(argumentString, " ")
}

func parseArguments(args []string, arguments []Argument) (map[string]Argument, error) {
	returnArguments := map[string]Argument{}
	if err := validateArguments(arguments); err != nil {
		return returnArguments, err
	}

	maxArgs := len(arguments)
	minArgs := 0
	for _, argument := range arguments {
		if !argument.Optional {
			minArgs++
		}
	}

	argumentWord := "argument"
	if maxArgs != 1 {
		argumentWord = "arguments"
	}
	errorMessage := fmt.Sprintf("This command requires %d", minArgs)
	if minArgs != maxArgs {
		errorMessage = fmt.Sprintf("%s and at most %d %s", errorMessage, maxArgs, argumentWord)
	}

	errorMessage = fmt.Sprintf("%s: %s", errorMessage, argumentString(arguments))

	if len(args) == 0 {
		if len(arguments) == 0 {
			return returnArguments, nil
		}

		if !arguments[0].Optional {
			return returnArguments, fmt.Errorf(errorMessage)
		}
	}

	for i, value := range args {
		arguments[i].HasValue = true
		arguments[i].Value = value
	}

	for _, argument := range arguments {
		if argument.Value == nil {
			if argument.Type == ArgumentBool {
				argument.Value = false
			} else if argument.Type == ArgumentInt {
				argument.Value = 0
			} else if argument.Type == ArgumentString {
				argument.Value = ""
			}
			argument.HasValue = false
		}
		returnArguments[argument.Name] = argument
	}

	return returnArguments, nil
}

func validateArguments(arguments []Argument) error {
	reachedOptional := false
	for _, arg := range arguments {
		if reachedOptional {
			if !arg.Optional {
				return fmt.Errorf("Argument %s must be placed before all optional arguments", arg.Name)
			}
			continue
		}

		if arg.Optional {
			reachedOptional = true
		}
	}
	return nil
}
