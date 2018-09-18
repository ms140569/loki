package level

import (
	"testing"
)

func TestLevelNumbers(t *testing.T) {
	if len(levels) != 8 {
		t.Fatal("Number of Loglevels changed.")
	}
}

func TestSomeLevels(t *testing.T) {
	if ToLoglevel("off") != Off {
		t.Fatal("Loglevel Off not recognized.")
	}

	if ToLoglevel("Trace") != Trace {
		t.Fatal("Loglevel Off not recognized.")
	}
}
