package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"loki/config"
	"loki/log"
	"loki/record"
	pb "loki/storage"
	"loki/subcommand"
	"loki/utils"
	"os"
	"strings"
	"time"
)

// Import lets you import Keepass CSV export files into the Loki store.
// Example:
// loki import ~/prj/golang/src/loki/data/import/keepass-export.csv
func Import(cfg config.Configuration, subcommand subcommand.Subcommand, args ...string) error {

	if len(args) > 1 {
		return errors.New("Too many arguments given")
	}

	filename := args[0]

	log.Info("File to import: " + filename)

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	/*
	   0 -> Group
	   1 -> Title
	   2 -> Username
	   3 -> Password
	   4 -> URL
	   5 -> Notes
	*/

	var records = make(map[string]pb.Record)

	r := csv.NewReader(strings.NewReader(string(data)))

	skipped := false

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("%v", err)
		}

		// Skip the first line

		if !skipped {
			skipped = true
			continue
		}

		group := record[0]
		title := record[1]
		username := record[2]
		password := record[3]
		url := record[4]
		notes := record[5]

		log.Info("--------------------------------")
		log.Info("Group    -> " + group)
		log.Info("Title    -> " + title)
		log.Info("Username -> " + username)
		log.Info("Password -> " + password)
		log.Info("URL      -> " + url)
		log.Info("Notes    -> " + notes)

		filename := groupAndTitleToFilename(group, title)

		// create one tags entry with the date of today to indicate this
		// entry to be created by import
		t := time.Now()
		nowAsString := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
		tags := make([]string, 0)
		tags = append(tags, "imported-"+nowAsString)

		records[filename] = pb.Record{Title: title, Account: username, Password: password, Tags: tags, Url: url, Notes: notes}
	}

	cnt := len(records)

	if cnt > 0 {
		log.Info("\nImporting %d records\n", cnt)
		key, _ := utils.GetMasterkey(true)

		for filename, rec := range records {
			utils.CreateLeadingDirectories(filename)
			err := record.WriteRecord(filename, cfg.Generation, key, rec)

			if err != nil {
				return err
			}
		}
		utils.SetupKeyAgent(key)
	}

	return nil
}

func groupAndTitleToFilename(group string, title string) string {
	if group == "Root" {
		return utils.NormalizePath(title)
	}
	dir := strings.TrimPrefix(group, "Root/")
	return utils.NormalizePath(dir + string(os.PathSeparator) + title)
}
