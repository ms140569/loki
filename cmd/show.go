package cmd

import (
	"loki/config"
	"loki/log"
	"loki/record"
	"loki/subcommand"
	"loki/utils"
	"github.com/atotto/clipboard"
	"os"
)

// Show displays a single record (lokifile) with all its content.
func Show(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	filename := utils.NormalizePath(args[0])

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Error("File does not exist: %v", err)
		return err
	}

	key, _ := utils.GetMasterkey(false)

	rec, hdr, err := record.LoadRecord(filename, key)

	if err != nil {
		log.Error("Error reading record: %v", err)
		return err
	}

	log.Info("Record: %s\n", filename)
	hdr.Print(0)

	utils.Display(rec, cfg.Blindmode)

	if cfg.Clipboard {
		clipboard.WriteAll(rec.Password)
	}

	utils.SetupKeyAgent(key)

	return nil
}
