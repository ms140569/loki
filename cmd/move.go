package cmd

import (
	"loki/config"
	"loki/log"
	"loki/subcommand"
	"loki/utils"
	"fmt"
	"os"
	"strings"
)

// Move moves loki password file from one location to another.
// file1 -> file2
// file1 -> dir/file1
// file1 -> dir/file2
// dir1 -> dir2
// dir1 -> file1 ***** ERROR *****
func Move(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	key, _ := utils.GetMasterkey(false)

	if key == nil {
		return fmt.Errorf("could not get Masterkey")
	}

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
			return fmt.Errorf("Destination should be no filename when moving directories: %s -> %s", src, dst)
		}
	} else {

		info, err = os.Stat(args[1])

		if err == nil && info.IsDir() {
			dst = args[1] + string(os.PathSeparator) + src
		} else {
			dst = utils.NormalizePath(args[1])
		}
	}

	log.Info("Moving: %s -> %s", src, dst)

	err = os.Rename(src, dst)

	if err != nil {
		return err
	}

	utils.SetupKeyAgent(key)
	return nil
}
