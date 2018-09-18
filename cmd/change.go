package cmd

import (
	"errors"
	"fmt"
	"loki/config"
	"loki/log"
	"loki/record"
	"loki/subcommand"
	"loki/tree"
	"loki/utils"
)

// ChangeMasterkey changes the password for all files in the store. This modifies all every single file plus the .master file.
func ChangeMasterkey(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {

	if len(args) > 0 {
		return errors.New("Too many arguments given")
	}

	log.Info("Change Masterkey.")

	base := cfg.SystemDirectory()
	oldgeneration := cfg.Generation

	if !utils.CheckBase(cfg) {
		return fmt.Errorf("Problem getting basedir")
	}

	// Stop Agent
	utils.ShutdownAgent()

	// Verify old key

	log.Info("Please provide the old password for verification.")

	oldkey, err := utils.GetMasterkeyWithAgent(false, false)

	if err != nil {
		return err
	}

	// generate map of all files:

	fm := tree.CreateFilemap(base, oldkey)

	items := len(*fm)

	if items < 1 {
		log.Info("No data to change. Exit.")
		return nil
	}

	log.Info("Found %d items to change:\n", items)

	// Request new key twice
	log.Info("Please provide the NEW password.")
	newkey, _ := utils.GetMasterkeyWithAgent(true, false)

	// Changes all files
	for k, v := range *fm {
		log.Debug("key[%s] value[%s]\n", k, v)
		record.WriteRecord(k, oldgeneration+1, newkey, *v)
	}

	// Update Masterfile to indicate global change
	utils.RaiseGenerationInMasterfile(cfg.GetMasterfilename())

	// Verify all files
	if err := tree.Verify(base, newkey); err != nil {
		log.Error("Tree verification failed: %v", err)
	}

	return nil
}
