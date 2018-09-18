package cmd

import (
	"fmt"
	"loki/config"
	"loki/log"
	"loki/record"
	pb "loki/storage"
	"loki/subcommand"
	"loki/utils"
	"os"
	"reflect"
)

// Diff diffs two loki files. To be used in conjunction with a vcs like git.
func Diff(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {
	oldpath := utils.NormalizePath(args[0])
	newpath := utils.NormalizePath(args[1])

	for _, filename := range []string{oldpath, newpath} {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Error("File does not exist: %v", err)
			return err
		}
	}

	key, err := utils.GetMasterkey(false)

	if err != nil {
		return fmt.Errorf("no acces")
	}

	diff, err := diffFiles(oldpath, newpath, key)

	if err != nil {
		return err
	}

	log.Info(diff)

	utils.SetupKeyAgent(key)
	return nil
}

func diffFiles(oldpath, newpath string, key []byte) (string, error) {

	var old *pb.Record
	var new *pb.Record
	var err error

	old, _, err = record.LoadRecord(oldpath, key)

	if err != nil {
		log.Error("Error reading record: %v", err)
		return "", err
	}

	new, _, err = record.LoadRecord(newpath, key)

	if err != nil {
		log.Error("Error reading record: %v", err)
		return "", err
	}

	return diffRecords(old, new), nil
}

func diffRecords(old, new *pb.Record) string {

	type fieldDesc struct {
		name  string
		label string
	}

	var diff string

	fields := map[int]fieldDesc{
		0: fieldDesc{"Title", config.TitleLabel},
		1: fieldDesc{"Account", config.AccountLabel},
		2: fieldDesc{"Password", config.PasswordLabel},
		3: fieldDesc{"Url", config.URLLabel},
		4: fieldDesc{"Notes", config.NotesLabel},
	}

	ov := reflect.ValueOf(*old)
	nv := reflect.ValueOf(*new)

	for idx := range []int{0, 1, 2, 3, 4} {

		field := fields[idx].name
		label := fields[idx].label

		ovalue := fmt.Sprintf("%s", ov.FieldByName(field))
		nvalue := fmt.Sprintf("%s", nv.FieldByName(field))

		if ovalue != nvalue {
			diff = diff + fmt.Sprintf("-"+label+"%s\n", ovalue)
			diff = diff + fmt.Sprintf("+"+label+"%s\n", nvalue)
		}
	}
	return diff
}
