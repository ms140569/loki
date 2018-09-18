package cmd

import (
	"loki/config"
	"loki/log"
	"loki/record"
	"loki/subcommand"
	"loki/tree"
	"loki/utils"
	"errors"
	"os"
	"strings"
)

// Dump dumps the full content of all loki-files in the tree to the console for examination. This is mainly
// for debugging purposes, therefore this is a hidden command.
func Dump(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {

	base := cfg.SystemDirectory()

	if !utils.CheckBase(cfg) {
		return errors.New("could not find basedir")
	}

	key := dumpWalker(base, cfg.Blindmode)

	utils.SetupKeyAgent(key)

	return nil
}

func dumpWalker(dir string, blind bool) []byte {

	var key []byte

	tree.FilteredWalk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			relPath := strings.TrimPrefix(path, dir)

			// count slashes
			column := (strings.Count(relPath, string(os.PathSeparator)) - 1) * 4

			log.Info("\n\n")
			log.Info("------------------------------------------------------------------------------")
			log.Info("Path: " + relPath)
			log.Info("------------------------------------------------------------------------------")

			if key == nil {
				key, _ = utils.GetMasterkey(false)
			}

			rec, hdr, err := record.LoadRecord(path, key)

			if err != nil {
				log.Error("Error reading record: %v", err)
				return nil
			}

			log.Info("\n%*s---- Header ------\n", column, "")
			hdr.Print(column)
			log.Info("")

			utils.PrefixedDisplay(rec, column, blind)

		}

		return nil
	})

	return key
}
