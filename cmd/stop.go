package cmd

import (
	"errors"
	"loki/config"
	"loki/log"
	"loki/subcommand"
	"loki/utils"
)

// Stop stops the loki-agent running in the background.
func Stop(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {

	if len(args) > 0 {
		return errors.New("Too many arguments given")
	}

	log.Info("Stopping agent")
	utils.ShutdownAgent()
	return nil
}
