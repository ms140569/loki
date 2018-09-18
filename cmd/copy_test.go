package cmd

import (
	"loki/utils"
	"testing"
)

func TestCopyCase1_file_file(t *testing.T) {
	defer SetupTest(t)()
	if err := Copy(cfg, cmd, "file1", "file3"); err != nil {
		t.Fail()
	}

	if !utils.VerifyFile(TBASE() + "file3.loki") {
		t.Fail()
	}
}

func TestCopyCase2_file_dir(t *testing.T) {
	defer SetupTest(t)()
	if err := Copy(cfg, cmd, "file1", "dir1"); err != nil {
		t.Fail()
	}

	if !utils.VerifyFile(TBASE() + "dir1" + SEP + "file1.loki") {
		t.Fail()
	}
}

func TestCopyCase3_dir_file(t *testing.T) {
	defer SetupTest(t)()
	if err := Copy(cfg, cmd, "dir1", "file1"); err == nil {
		t.Fail()
	}
}

func TestCopyCase4_dir_dir(t *testing.T) {
	defer SetupTest(t)()
	if err := Copy(cfg, cmd, "dir1", "dir4"); err != nil {
		t.Fail()
	}

	if !utils.VerifyFile(TBASE() + "dir4" + SEP + "frumpy.loki") {
		t.Fail()
	}
}

/*
func TestCopyCase4_dir_dir_exists(t *testing.T) {
	if err := Copy(cfg, cmd, "dir1", "dir3"); err != nil {
		t.Fail()
	}

	if !utils.VerifyFile(tmpDir + string(os.PathSeparator) + "dir3" + string(os.PathSeparator) + "dir1" + string(os.PathSeparator) + "frumpy.loki") {
		t.Fail()
	}
}
*/
