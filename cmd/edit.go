package cmd

import (
	"loki/config"
	"loki/log"
	"loki/record"
	"loki/subcommand"
	"loki/utils"
)

// Edit lets you edit a single record ( lokifile ). You might use an external editor for the
// Notes field by setting the EDITOR environment variable.
func Edit(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	filename := utils.NormalizePath(args[0])

	key, _ := utils.GetMasterkey(false)

	rec, hdr, err := record.LoadRecord(filename, key)

	if err != nil {
		log.Error("Error reading record: %v", err)
		return err
	}

	hdr.Print(0)

	err = rec.Edit(!cfg.ExternalEditor)

	if err != nil {
		return err
	}

	if cfg.ExternalEditor {
		notes, err := utils.StartEditorWithData(rec.Notes)

		if err != nil {
			log.Error("Error starting external editor: %v", err)
			return err
		}

		rec.Notes = notes
	}

	record.WriteRecord(filename, cfg.Generation, key, *rec)

	utils.SetupKeyAgent(key)

	return nil
}
