package cmd

import (
	"loki/config"
	"loki/log"
	"loki/subcommand"
	"loki/utils"
)

// Init initializes a new Loki password-manager directory. This is usually ~/.loki.
// or it is given as an argument.
func Init(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	if len(args) > 0 && len(args[0]) > 0 {
		cfg.SetSystemDirectory(args[0])
	}

	err := utils.InitBasedir(cfg)
	basedir := cfg.SystemDirectory()

	if err != nil {
		log.Error("Problem initializing basedir: %s, error: %v", basedir, err)
		return err
	}

	log.Info("Sucessfully initialized basedir : %s", basedir)
	return nil
}
