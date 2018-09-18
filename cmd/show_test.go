package cmd

import (
	"testing"
)

func TestCommandShow(t *testing.T) {
	defer SetupTest(t)()
	err := Show(cfg, cmd, "file1")

	if err != nil {
		t.Fail()
	}
}

func TestCommandShow2(t *testing.T) {
	defer SetupTest(t)()
	err := Show(cfg, cmd, "dir3/sub/bingo.loki")

	if err != nil {
		t.Fail()
	}
}
