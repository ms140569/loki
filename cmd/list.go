package cmd

import (
	"errors"
	"github.com/xlab/treeprint"
	"loki/config"
	"loki/log"
	"loki/subcommand"
	tu "loki/tree"
	"loki/utils"
	"os"
	"path/filepath"
	"strings"
)

type treeMap map[string]treeprint.Tree

// List displays the contents of the password store in a treelike fashion.
func List(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {

	base := cfg.SystemDirectory()

	if !utils.CheckBase(cfg) {
		return errors.New("could not find basedir")
	}

	// possibly adding a subdir

	if len(args) > 0 && utils.VerifyDirectory(base+string(os.PathSeparator)+args[0]) {
		base += string(filepath.Separator) + args[0]
	}

	walker(base)

	return nil
}

func walker(dir string) {

	tm := make(treeMap)
	tree := treeprint.New()

	tu.FilteredWalk(dir, func(path string, info os.FileInfo, err error) error {
		relPath := strings.TrimPrefix(path, dir)

		parent := lookupParent(tm, relPath)

		if parent == nil {
			parent = tree
		}

		if len(relPath) == 0 {
			return nil
		}

		if info.IsDir() {
			tm[relPath] = parent.AddBranch(info.Name())
		} else {
			parent.AddNode(strings.TrimSuffix(info.Name(), config.FileSuffix))
		}

		return nil
	})

	output := tree.String()

	if output == ".\n" {
		log.Info("no data.")
	} else {
		log.Info(output)
	}
}

func lookupParent(tm treeMap, path string) treeprint.Tree {
	if len(path) == 0 {
		return nil
	}

	dir, _ := filepath.Split(path)

	dir = strings.TrimSuffix(dir, string(os.PathSeparator))

	for key, val := range tm {
		log.Debug("Dir: %s, Key: %s\n", dir, key)
		if key == dir {
			return val
		}
	}
	return nil
}
