package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"loki/log"
)

// https://github.com/otiai10/copy

// Copy copies either files or directories.
// FileInfo for the src must be given to achieve this.
func Copy(src, dst string, info os.FileInfo) error {
	if info.IsDir() {
		log.Debug("Directory copy (%s -> %s)", src, dst)
		return dcopy(src, dst, info)
	}
	log.Debug("Filecopy (%s -> %s)", src, dst)
	return fcopy(src, dst, info)
}

func fcopy(src, dest string, info os.FileInfo) error {

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

func dcopy(src, dest string, info os.FileInfo) error {

	// if destination is a directory, we gotta copy *into* it.
	/*
		i, e := os.Stat(dest)

		if e == nil && i.IsDir() && info.IsDir() {
			dest = filepath.Join(dest, src)
			log.Debug("Copy into: %s", dest)
		}
	*/

	if err := os.MkdirAll(dest, info.Mode()); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if err := Copy(
			filepath.Join(src, info.Name()),
			filepath.Join(dest, info.Name()),
			info,
		); err != nil {
			return err
		}
	}

	return nil
}
