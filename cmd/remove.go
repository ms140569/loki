package cmd

import (
	"loki/config"
	"loki/log"
	"loki/subcommand"
	"loki/utils"
	"errors"
	"os"
)

// Remove removes either a single password file or a subtree from the store.
func Remove(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	filename := args[0]

	if isDir(filename) {
		key, _ := utils.GetMasterkey(false)

		if key == nil {
			return errors.New("could not get Masterkey")
		}

		err := os.RemoveAll(filename)

		if err != nil {
			log.Error("Error removing directory: %v", err)
			return err
		}
		utils.SetupKeyAgent(key)
		return nil
	}

	filename = utils.NormalizePath(filename)

	_, err := os.Stat(filename)

	if err != nil {
		log.Error("Path does not exist : " + filename)
		return err
	}

	key, _ := utils.GetMasterkey(false)

	if key == nil {
		return errors.New("could not get Masterkey")
	}

	err = os.Remove(filename)

	if err != nil {
		return err

	}

	utils.SetupKeyAgent(key)
	return nil
}

func isDir(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fi.IsDir()
}
