package cmd

import (
	"loki/log"
	"loki/utils"
	"os"
	"testing"
)

func TestCommandInit(t *testing.T) {
	defer SetupTest(t)()
	tmpDir, err := utils.CreateTempdirWithPrefx("init-test")

	if err != nil {
		log.Error("%v", err)
		return
	}

	log.Debug("Init here: %s", tmpDir)

	cfg.SetSystemDirectory(TBASE() + "loki")

	err = Init(cfg, cmd)

	if err != nil {
		log.Error("%v", err)
		t.Fail()
	}

	if !utils.VerifyFile(TBASE() + "loki" + SEP + ".master") {
		t.Fail()
	}

	if !utils.VerifyFile(TBASE() + "loki" + SEP + ".config") {
		t.Fail()
	}

	if err = os.RemoveAll(tmpDir); err != nil {
		log.Error("%v", err)
	}

}
