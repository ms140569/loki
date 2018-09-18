package cmd

import (
	"loki/utils"
	"testing"
)

func TestRemove1_file(t *testing.T) {
	defer SetupTest(t)()
	if err := Remove(cfg, cmd, "file1"); err != nil {
		t.Fail()
	}

	if utils.VerifyFile(TBASE() + "file1.loki") {
		t.Fail()
	}
}

func TestRemove2_dir(t *testing.T) {
	defer SetupTest(t)()
	if err := Remove(cfg, cmd, "dir1"); err != nil {
		t.Fail()
	}

	if utils.VerifyFile(TBASE() + "dir1") {
		t.Fail()
	}
}
