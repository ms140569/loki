package cmd

import (
	"loki/config"
	"loki/log"
	"loki/record"
	"loki/subcommand"
	"loki/tree"
	"loki/utils"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Search searches a term case-insensitive in all fields and the pathname. If found the string gets highligted.
func Search(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	searchstring := strings.ToLower(args[0])

	base := cfg.SystemDirectory()

	if !utils.CheckBase(cfg) {
		return errors.New("could not find basedir")
	}

	var key []byte
	var err error

	if key, err = utils.GetMasterkey(false); err != nil {
		return fmt.Errorf("Problem getting masterkey: %v", err)
	}

	if searchWalker(base, key, searchstring, cfg.Blindmode) == nil {
		utils.SetupKeyAgent(key)
	}

	return nil
}

func searchWalker(dir string, key []byte, searchstring string, blind bool) error {

	var outError error

	tree.FilteredWalk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {

			relPath := strings.TrimPrefix(path, dir+string(os.PathSeparator))

			rec, _, err := record.LoadRecord(path, key)

			if err != nil {
				log.Error("Error reading record. File: %s, Error: %v", path, err)
				outError = err
				return io.EOF
			}

			if rec.Search(searchstring) || strings.Contains(strings.ToLower(relPath), searchstring) {
				log.Info("Record: %s\n", utils.Highlight(relPath, searchstring))

				log.Debug("Found searchstring : " + searchstring)
				utils.PrefixedDisplayWithHighlighting(rec, 0, searchstring, blind)
				log.Info("")
			}

		}

		return nil
	})

	return outError
}
