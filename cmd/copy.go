package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"loki/config"
	"loki/log"
	"loki/subcommand"
	"loki/utils"
)

// Copy copies a single pasword (*.loki) to a new location. Directory copies are not supported yet.
func Copy(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	key, _ := utils.GetMasterkey(false)

	if key == nil {
		return errors.New("could not get Masterkey")
	}

	// source is either:
	// - wrong
	// - a directory
	// - a file ( + .loki )

	srcIsDir := false

	src := args[0]

	// 1st try on source
	info, err := os.Stat(src)

	if err != nil {
		log.Debug("Source is either wrong or an un-normalized file: %s", src)
		src = utils.NormalizePath(args[0])

		// 2nd try on source
		info, err = os.Stat(src)

		if err != nil {
			return fmt.Errorf("Could not find source: %s[.loki]", args[0])
		}
	}

	// src is either a dir or the full filename including *.loki was given
	if info.IsDir() {
		srcIsDir = true
	}

	dst := args[1]

	if srcIsDir { // ... dst should be a directory as well
		if strings.HasSuffix(dst, config.FileSuffix) || utils.VerifyFile(utils.NormalizePath(dst)) {
			return fmt.Errorf("Destination should be no filename when copying directories: %s -> %s", src, dst)
		}
	} else {
		// if src is a file the dst might be:
		// - a directory
		// - a file prefixed with a directory

		dstInfo, dstErr := os.Stat(dst)

		if !(dstErr == nil && dstInfo.IsDir()) {
			dst = utils.NormalizePath(args[1])

			if strings.Contains(dst, "/") {
				utils.CreateLeadingDirectories(dst)
			}
		} else {
			// dst is an directory. Adding full destination filename
			dst = dst + string(os.PathSeparator) + filepath.Base(src)
		}
	}

	log.Info("Copying: %s -> %s", src, dst)
	err = utils.Copy(src, dst, info)

	if err != nil {
		return fmt.Errorf("Error copying: %v", err)
	}

	utils.SetupKeyAgent(key)
	return nil
}
