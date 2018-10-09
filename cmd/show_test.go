package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestCommandShow(t *testing.T) {
	t.Skip("skipping test for now, further investigation needed")
	defer SetupTest(t)()
	err := Show(cfg, cmd, "file1")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error showing file: %v", err)
		t.Fail()
	}
}

func TestCommandShow2(t *testing.T) {
	t.Skip("skipping test for now, further investigation needed")
	defer SetupTest(t)()
	err := Show(cfg, cmd, "dir3/sub/bingo.loki")

	if err != nil {
		fmt.Fprintf(os.Stderr, "error showing file: %v", err)
		t.Fail()
	}
}
