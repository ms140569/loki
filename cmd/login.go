package cmd

import (
	"loki/config"
	"loki/crypto"
	"loki/log"
	"loki/subcommand"
	"loki/tree"
	"loki/utils"
)

var kdf = crypto.NewKeyDerivator()

// Login lets you verify password against the loki-store and thereby starting an key-agent for your convienience.
func Login(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	key, _ := utils.GetMasterkey(false)
	base := cfg.SystemDirectory()

	rec := tree.GetFirstRecord(base, key)

	if rec == nil {
		log.Info("No data found, therefore no login.")
		return nil
	}

	utils.SetupKeyAgent(key)
	return nil
}
