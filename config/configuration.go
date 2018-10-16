package config

import (
	"flag"
	"fmt"
	"gopkg.in/gcfg.v1"
	"loki/log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Custom boolean Flag-type to distinguish between default values
// and no flag given at all, see:
// https://stackoverflow.com/questions/35809252/check-if-flag-was-provided-in-go

type boolFlag struct {
	set   bool
	value bool
}

func (bf *boolFlag) Set(x string) error {

	if strings.ToLower(x) == "true" {
		bf.value = true
	} else {
		bf.value = false
	}

	bf.set = true
	return nil
}

func (bf *boolFlag) IsBoolFlag() bool {
	return true
}

func (bf *boolFlag) IsSet() bool {
	return bf.set
}

func (bf *boolFlag) String() string {
	if bf.value {
		return "true"
	}
	return "false"
}

// FlagBundle is to bundle all flags the programm understands into
// one entity to be passed around together.
type FlagBundle struct {
	Loglevel       string
	ExternalEditor boolFlag
	Gitmode        boolFlag
	Clipboard      boolFlag
	Blindmode      boolFlag
	Debug          boolFlag
	Help           boolFlag
}

// Configuration is the system configuration as created by merging config-file values and
// commandline-flags passed around with the FlagBundle.
type Configuration struct {
	SystemDir      string
	Binpath        string
	Gitmode        bool
	Generation     uint32
	Loglevel       string
	ExternalEditor bool
	Clipboard      bool
	Blindmode      bool
}

// ParseFlags defines all flags the program understands, parses the commandline into them and
// returns them as a bundle of flags
func ParseFlags() FlagBundle {
	fb := FlagBundle{}

	flag.StringVar(&fb.Loglevel, "l", "INFO", "Loglevel the program is running with.")

	flag.Var(&fb.ExternalEditor, "e", "Use external editor given in the EDITOR environment variable.")
	flag.Var(&fb.Gitmode, "g", "Automatically run git commit after each modifiying command.")
	flag.Var(&fb.Clipboard, "c", "Copy password to clipboard")
	flag.Var(&fb.Blindmode, "b", "Blindmode. Do not show password.")
	flag.Var(&fb.Debug, "d", "Debug mode. Equivalent to -l debug.")
	flag.Var(&fb.Help, "h", "Show help information.")

	flag.Parse()

	return fb
}

// New creates a new system-configuration *AND* parses the os.Args vector into a flagbundle.
func New(fb FlagBundle) Configuration {

	var loglevel string

	if fb.Debug.set {
		loglevel = "Debug"
	} else {
		loglevel = fb.Loglevel
	}

	configFilename := getSystemDirectory() + string(os.PathSeparator) + ConfigFilename

	_, err := os.Stat(configFilename)

	var cfg Configuration

	if err != nil {
		// log.Debug("No configfile found.")
		cfg = Configuration{}
	} else {
		// log.Debug("Found configfile : %s", configFilename)
		cfg = readConfigFile(configFilename)

	}

	cfg.Loglevel = loglevel

	// Flags override configfile-values *if* they are provided:

	if fb.ExternalEditor.set {
		cfg.ExternalEditor = fb.ExternalEditor.value
	}

	if fb.Gitmode.set {
		cfg.Gitmode = fb.Gitmode.value
	}

	if fb.Clipboard.set {
		cfg.Clipboard = fb.Clipboard.value
	}

	if fb.Blindmode.set {
		cfg.Blindmode = fb.Blindmode.value
	}

	return cfg
}

// Print prints the system-configuration to the console with DEBUG-level
func (c *Configuration) Print() {
	log.Debug("Binpath    : %s", c.Binpath)
	log.Debug("Generation : %d", c.Generation)

	log.Debug("Gitmode    : %t", c.Gitmode)
	log.Debug("Clipboard  : %t", c.Clipboard)
	log.Debug("ExtEditor  : %t", c.ExternalEditor)
	log.Debug("Loglevel   : %s\n", c.Loglevel)
}

// GreetingString returns the softwares greeeting string (including the version).
func (c *Configuration) GreetingString() string {
	return "Loki Password Manager, ver " + SoftwareVersion
}

// SystemDirectory returns the password store directory this system is working with.
func (c *Configuration) SystemDirectory() string {
	if len(c.SystemDir) > 0 {
		return c.SystemDir
	}
	return getSystemDirectory()
}

// SetSystemDirectory is used for testing to set an individual testdirectory.
func (c *Configuration) SetSystemDirectory(dirname string) {
	c.SystemDir = dirname
}

// GetMasterfilename returns the full path to the systems masterfile. Usually: ~/.loki/.master.
func (c *Configuration) GetMasterfilename() string {
	if len(c.SystemDir) > 0 {
		return c.SystemDir + string(os.PathSeparator) + MasterFilename
	}
	return getSystemDirectory() + string(os.PathSeparator) + MasterFilename
}

func GetSocketfilePath() string {
	return fmt.Sprintf(CommunicationFile, os.Getuid())
}

func getSystemDirectory() string {
	baseVar := os.Getenv(LokiBaseEnv)

	if len(baseVar) > 0 {
		return baseVar
	}
	usr, _ := user.Current()
	dir := usr.HomeDir

	return filepath.Join(dir, DefaultDirectory)
}

func readConfigFile(filename string) Configuration {

	// We have to give a struct-in-struct here to match the ini-style sections
	type Config struct {
		Basic Configuration
	}

	cfg := Config{}
	err := gcfg.ReadFileInto(&cfg, filename)

	if err != nil {
		log.Error("Error reading file: %s", err.Error())
	}

	return cfg.Basic
}
