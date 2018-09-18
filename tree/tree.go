package tree

import (
	"loki/log"
	"loki/record"
	pb "loki/storage"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileMap is the basic mapping structure
// filepath -> Record to be used all over the place
type FileMap map[string]*pb.Record

// FilteredWalk filteres hidden files and the .git directory from a filewalker. This is used
// by the subcommands which walk the whole tree: list, change, dump and search.
func FilteredWalk(dir string, walkFn filepath.WalkFunc) error {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		relPath := strings.TrimPrefix(path, dir)
		log.Debug("Relpath: \"%s\", dir: %t, Name: %s, Syspath: %s\n", relPath, info.IsDir(), info.Name(), path)

		if strings.HasPrefix(relPath, ".git") || strings.HasPrefix(relPath, "/.git") || strings.HasPrefix(info.Name(), ".") {
			return nil
		}
		return walkFn(path, info, err)
	})
	return nil
}

// CreateFilemap produces a map of filenames -> records of all loki-files in the
// datastore.
func CreateFilemap(dir string, key []byte) *FileMap {

	fm := make(FileMap)

	FilteredWalk(dir, func(path string, info os.FileInfo, err error) error {
		relPath := strings.TrimPrefix(path, dir)

		if len(relPath) == 0 {
			return nil
		}

		if !info.IsDir() {
			rec, _, err := record.LoadRecord(path, key)

			if err != nil {
				return errors.New("could not load record")
			}
			fm[path] = rec
		}

		return nil
	})

	return &fm
}

// GetFirstRecord walks the given directory tree and returns the
// very first record found or nil.
func GetFirstRecord(dir string, key []byte) *pb.Record {

	var first *pb.Record

	FilteredWalk(dir, func(path string, info os.FileInfo, err error) error {
		relPath := strings.TrimPrefix(path, dir)

		if len(relPath) == 0 {
			return nil
		}

		if !info.IsDir() {
			rec, _, err := record.LoadRecord(path, key)

			if err != nil {
				return errors.New("could not load record")
			}

			first = rec
			return io.EOF
		}

		return nil
	})

	return first
}

// Verify create a filemap of all records in the tree given by base and
// and unlocked by the parameter key. This alone should verify the tree
// but to be save we show the title in addition.
func Verify(base string, key []byte) error {

	fm := CreateFilemap(base, key)

	// Changes all files
	for k, rec := range *fm {
		log.Debug("Verifying entry: %s -> %s", k, rec.Title)
	}

	log.Debug("Verify %d records.", len(*fm))
	return nil
}
