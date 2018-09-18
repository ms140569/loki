package cmd

import (
	"loki/log"
	"loki/utils"
	"testing"
)

func TestMoveCase1_file_file(t *testing.T) {
	defer SetupTest(t)()
	if err := Move(cfg, cmd, "file2", "file3"); err != nil {
		t.Fail()
	}

	if !utils.VerifyFile(TBASE() + "file3.loki") {
		t.Fail()
	}
}

func TestMoveCase2_file_dir(t *testing.T) {
	defer SetupTest(t)()
	if err := Move(cfg, cmd, "file2", "dir2"); err != nil {
		t.Fail()
	}

	if !utils.VerifyFile(TBASE() + "dir2" + SEP + "file2.loki") {
		t.Fail()
	}
}

func TestMoveCase3_dir_file_negative(t *testing.T) {
	defer SetupTest(t)()
	if err := Move(cfg, cmd, "dir2", "file1"); err == nil {
		t.Fail()
	}
}

func TestMoveCase4_dir_dir(t *testing.T) {
	defer SetupTest(t)()
	if err := Move(cfg, cmd, "dir2", "dir5"); err != nil {
		log.Debug("Move returned error.")
		log.Error("Error: %v", err)
		t.Fail()
	}

	if !utils.VerifyDirectory(TBASE() + "dir5") {
		log.Debug("Ended here?")
		t.Fail()
	}
	log.Debug("All is fine.")
}
