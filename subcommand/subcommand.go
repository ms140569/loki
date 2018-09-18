package subcommand

import (
	"errors"
	"github.com/fatih/color"
	"loki/config"
	"loki/log"
	"strings"
)

// Handler is a type of function which could be used to process a registered loki password manager subcommand
type Handler func(config.Configuration, Subcommand, ...string) error

// Subcommand is the records which holds all information to register, call and display loki password manager subcommands
type Subcommand struct {
	Aliases        []string // Aliases a command is available at. Example: loki list | ls. loki version | ver
	NumberOfParams int      // minimal number of parameter a command needs. Could be verified easily.
	Paramhint      string   // hint what a prameter is: file, dir, ...
	DefaultCommand bool     // is this the command to run, if *NO* command is given? Should only be one. Only first one is used in lookup.
	Handler        Handler  // method to call if we dispatch this command.
	Description    string   // An explanation of the command which is used in the help function.
	Hidden         bool     // is this a hidden command? Excluded from help if true.
	Modifying      bool     // is this a store-modifying commmd? If yes it triggers a git commit in Gitmode.
}

func (subcommand Subcommand) help() {

	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgHiCyan).SprintFunc()

	log.Info(blue(strings.Join(subcommand.Aliases, " | ")) + " - " + subcommand.Description + "   Example: " + green(config.BinaryName) + " [flags] " + subcommand.Aliases[0] + " " + getParameterTemplate(subcommand.NumberOfParams, subcommand.Paramhint))

	log.Info("")
}

func (subcommand Subcommand) aliases() string {
	return strings.Join(subcommand.Aliases, " ") + " "
}

// CommandList is a registry of all registered subcommands known to the system.
type CommandList map[string]Subcommand

// Register makes given command with all its aliases known to the system. Using one of the aliases triggers a call to the handler.
func (cl CommandList) Register(aliases []string, numberOfParams int, paramHint string, defaultCommand bool, handler Handler, description string, hidden bool, modifying bool) {
	if len(aliases) == 0 {
		panic("Can not register without names")
	}
	// we register with the very first name
	basename := aliases[0]
	cl[basename] = Subcommand{aliases, numberOfParams, paramHint, defaultCommand, handler, description, hidden, modifying}
}

// PrintAll print a short help summary for all registered commands.
func (cl CommandList) PrintAll() {
	for _, v := range cl {
		if !v.Hidden {
			v.help()
		}
	}
}

// GetAllAliases generates a list of all aliases of all commands to be used in the bash programmable command substitiution.
func (cl CommandList) GetAllAliases() string {

	allaliases := ""

	for _, v := range cl {
		if !v.Hidden {
			allaliases += v.aliases()
		}
	}
	return allaliases
}

// FindCommand looks whether a command was registered with an alias given as name.
func (cl CommandList) FindCommand(name string) (Subcommand, error) {
	for _, v := range cl {
		for _, n := range v.Aliases {
			if n == name {
				return v, nil
			}
		}
	}
	return Subcommand{}, errors.New("Command not found :" + name)
}

// GetDefault returns the default command which should be run if the binary is executed without any parameter.
func (cl CommandList) GetDefault() (Subcommand, error) {
	for _, v := range cl {
		if v.DefaultCommand {
			return v, nil
		}
	}
	return Subcommand{}, errors.New("No default command given")
}

func getParameterTemplate(numberOfParams int, hint string) string {
	retVal := ""

	var paramHint string

	if hint != "" {
		paramHint = hint

		if numberOfParams == 0 {
			return hint
		}
	} else {
		paramHint = "param"
	}

	for i := 1; i <= numberOfParams; i++ {

		retVal += paramHint /* + strconv.Itoa(i) */ + " "
	}

	return retVal
}
