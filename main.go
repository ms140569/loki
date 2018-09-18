package main

import (
	"flag"
	"github.com/fatih/color"
	"loki/cmd"
	"loki/config"
	"loki/log"
	"loki/subcommand"
	"loki/utils"
	"os"
)

var commandList = make(subcommand.CommandList)

func main() {

	// Register all commands
	// Boolean: default command, hidden, modifying
	commandList.Register([]string{"ls", "list"}, 0, "", true, cmd.List, "Lists the password store in a treelike fashion.", false, false)
	commandList.Register([]string{"show"}, 1, "filename", false, cmd.Show, "Shows the contents of file.", false, false)
	commandList.Register([]string{"insert", "add"}, 1, "filename", false, cmd.Insert, "Inserts new data into file.", false, true)
	commandList.Register([]string{"import"}, 1, "keepass-filename", false, cmd.Import, "Imports a KeepassX CSV file.", false, true)
	commandList.Register([]string{"init"}, 0, "[pathname]", false, cmd.Init, "Initialize a new password store.", false, true)
	commandList.Register([]string{"login", "pw", "pass"}, 0, "", false, cmd.Login, "Authenticate against password store.", false, false)
	commandList.Register([]string{"dump"}, 0, "", false, cmd.Dump, "Dumps all information.", true, false)
	commandList.Register([]string{"search", "grep", "find"}, 1, "<querystring>", false, cmd.Search, "Searches for given string in all fields and recordnames.", false, false)
	commandList.Register([]string{"edit"}, 1, "filename", false, cmd.Edit, "Edit one Record.", false, true)
	commandList.Register([]string{"remove", "rm", "del"}, 1, "filename", false, cmd.Remove, "Delete a Record.", false, true)
	commandList.Register([]string{"copy", "cp"}, 2, "<file|dir>", false, cmd.Copy, "Copy a Record or a subtree.", false, true)
	commandList.Register([]string{"move", "mv"}, 2, "<file|dir>", false, cmd.Move, "Moves a Record or a subtree.", false, true)
	commandList.Register([]string{"shutdown", "stop"}, 0, "", false, cmd.Stop, "Stops the Agent.", false, false)
	commandList.Register([]string{"change"}, 0, "", false, cmd.ChangeMasterkey, "Changes the masterpassword in all files.", false, true)
	commandList.Register([]string{"diff"}, 2, "", false, cmd.Diff, "Diffs two files.", true, false)

	commandList.Register([]string{"help"}, 0, "", false, helpSubcommand, "Shows general help information.", false, false)
	// this will be transformed to a info command
	// commandList.Register([]string{"version", "ver"}, 0, "", false, versionSubcommand, "Shows version information.", false, false)
	// deactived for now, bash completion script shipped with debian-package. Might re-active in the future for systems with bash < 4
	// commandList.Register([]string{"complete"}, 0, "", false, completerSubcommand, "Generates bash programmable completion statement.", false, false)

	fb := config.ParseFlags()
	cfg := config.New(fb)
	log.SetSystemLogLevelFromString(cfg.Loglevel)

	arguments := flag.Args()

	sysdir := cfg.SystemDirectory()

	cfg.Binpath = utils.GetBinaryPath()

	log.Info("%s, data: %s\n", cfg.GreetingString(), sysdir)

	// display help in any case, wether we have a decent setup or not.

	if fb.Help.IsSet() || (len(arguments) == 1 && arguments[0] == "help") {
		help()
		utils.ExitSystemWithCode(0)
	}

	// Trying to change to the system directory

	if os.Chdir(sysdir) != nil {
		if len(arguments) > 0 && (arguments[0] == "init" || arguments[0] == "help") {
			utils.ExitSystem(dispatchToHandler(cfg, arguments[0], arguments[1:]...))
		}

		log.Fatal("Could not change to system directory.")
		log.Fatal("\nYou might want to initialize it with:")
		log.Fatal("  loki init [dirname]")

		utils.ExitSystemFailure()
	}

	// Load and print masterfile

	masterfile, err := utils.LoadMasterfile(cfg.GetMasterfilename())

	if err != nil {
		if len(arguments) < 1 || arguments[1] != "init" {
			log.Fatal("Could not read contents of masterfile, might do a init?")
			utils.ExitSystemFailure()
		}
	} else {
		cfg.Generation = masterfile.Generation
	}

	if err == nil {
		cfg.Print()
		log.Debug("\nMasterfile:\n\n")
		masterfile.Print(1)
	}

	if len(arguments) < 1 {
		dflt, err := commandList.GetDefault()

		if err != nil {
			log.Error("Could not find default command. Exit.")
			utils.ExitSystemFailure()
		}

		utils.ExitSystem(dispatchToHandler(cfg, dflt.Aliases[0]))
	} else {
		utils.ExitSystem(dispatchToHandler(cfg, arguments[0], arguments[1:]...))
	}
}

func dispatchToHandler(cfg config.Configuration, subcommand string, arg ...string) error {

	cmd, err := commandList.FindCommand(subcommand)

	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		log.Error(red("%s\n"), err.Error())
		help()
		return err
	}

	checkParams(cmd.NumberOfParams, len(arg))
	err = cmd.Handler(cfg, cmd, arg...)

	if cfg.Gitmode && cmd.Modifying && err == nil {

		// do a git add in any case. Should not hurt

		gitAdd := []string{"add", "."}
		log.Debug("Running git command: %v", gitAdd)
		utils.GitCommand(gitAdd)

		gitParams := []string{"commit", "-a", "-m", "'Loki commit triggered by command: " + cmd.Aliases[0] + "'"}
		log.Debug("Running git command: %v", gitParams)
		utils.GitCommand(gitParams)
	}

	return err
}

func helpSubcommand(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	help()
	return nil
}

func versionSubcommand(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	log.Info(cfg.GreetingString())
	return nil
}

func completerSubcommand(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	aliases := commandList.GetAllAliases()
	log.Info("\n\ncomplete -W \"%s\" loki\n\n", aliases)
	return nil
}

func help() {

	green := color.New(color.FgGreen).SprintFunc()
	log.Info("Run: %s [flags] command\n\n", green(config.BinaryName))
	commandList.PrintAll()

	log.Info("\nValid flags are:\n")

	flag.PrintDefaults()
	log.Info("")
}

func checkParams(needed int, given int) {
	if given < needed {
		log.Fatal("Not enough parameter given. Needed: %d, Given: %d\n\n", needed, given)
		help()
		utils.ExitSystemFailure()
	}
}
