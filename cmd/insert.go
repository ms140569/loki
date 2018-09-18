package cmd

import (
	"loki/config"
	"loki/log"
	"loki/record"
	"loki/storage"
	"loki/subcommand"
	"loki/utils"
	"errors"
)

// Insert adds a new record to the password store.
func Insert(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {

	if len(args) > 1 {
		return errors.New("Too many arguments given")
	}

	filename := utils.NormalizePath(args[0])

	file := utils.CreateLeadingDirectories(filename)

	log.Info("Filename: " + file)

	key, _ := utils.GetMasterkey(true)
	rec, err := storage.Ask()

	if err != nil {
		log.Info("Aborted.")
		return nil
	}

	err = record.WriteRecord(filename, cfg.Generation, key, rec)

	if err != nil {
		return err
	}

	utils.SetupKeyAgent(key)

	return nil
}
