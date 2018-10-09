package storage

import (
	"bufio"
	"fmt"
	"github.com/peterh/liner"
	"io"
	"loki/config"
	"loki/log"
	"os"
	"strings"
)

// Edit lets you edit the given record. Editing of the Notes field ends with CTRL-D.
func (rec *Record) Edit(editNotes bool) error {

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	pos := -1

	var err error
	var editedValue string

	aborted := fmt.Errorf("Aborted")

	if editedValue, err = line.PromptWithSuggestion(config.TitleLabel, rec.Title, pos); err != nil {
		return aborted
	}
	rec.Title = editedValue

	if editedValue, err = line.PromptWithSuggestion(config.AccountLabel, rec.Account, pos); err != nil {
		return aborted
	}
	rec.Account = editedValue

	if editedValue, err = line.PromptWithSuggestion(config.PasswordLabel, rec.Password, pos); err != nil {
		return aborted
	}
	rec.Password = editedValue

	tagsString := strings.Join(rec.Tags, ", ")

	if editedValue, err = line.PromptWithSuggestion(config.TagsLabel, tagsString, pos); err != nil {
		return aborted
	}
	rec.Tags = tagsStringToArray(editedValue)

	if editedValue, err = line.PromptWithSuggestion(config.URLLabel, rec.Url, pos); err != nil {
		return aborted
	}
	rec.Url = editedValue

	if editNotes {
		// handle multi-line notes
		log.Info("\n%s", config.NotesLabel)

		lines := strings.Split(rec.Notes, "\n")

		writeBack := ""
		lineIdx := 0

		aborted := false

		for {
			proposal := ""

			if lineIdx < len(lines) {
				proposal = lines[lineIdx]
			}

			input, err := line.PromptWithSuggestion("", proposal, len(proposal))

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Debug("Aborted? : %v", err)
				aborted = true
				break
			}

			writeBack += input + "\n"
			lineIdx++
		}

		if !aborted {
			rec.Notes = writeBack
		}
	}

	return nil
}

// Search searches case-insensitive in all fields of the record for the given string.
// It simply returns true if it's found anywhere.
func (rec *Record) Search(text string) bool {

	if strings.Contains(strings.ToLower(rec.Title), text) {
		return true
	}
	if strings.Contains(strings.ToLower(rec.Account), text) {
		return true
	}
	if strings.Contains(strings.ToLower(rec.Password), text) {
		return true
	}
	if strings.Contains(strings.ToLower(strings.Join(rec.Tags, ", ")), text) {
		return true
	}
	if strings.Contains(strings.ToLower(rec.Url), text) {
		return true
	}
	if strings.Contains(strings.ToLower(rec.Notes), text) {
		return true
	}

	return false
}

// split Tags like: "wlan, web, imported, mobile" into array
func tagsStringToArray(tagsString string) []string {
	tags := strings.Split(tagsString, ",")
	trimmedTags := make([]string, 0)

	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if len(trimmed) > 0 {
			log.Debug("Adding: |%s|", trimmed)
			trimmedTags = append(trimmedTags, trimmed)
		}
	}
	return trimmedTags
}

// Ask prompts the user for the content of all fields in the record.
func Ask() (Record, error) {
	rec := Record{}

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)

	var err error

	if rec.Title, err = line.Prompt(config.TitleLabel); err != nil {
		return rec, err
	}

	if rec.Account, err = line.Prompt(config.AccountLabel); err != nil {
		return rec, err
	}

	if rec.Password, err = line.Prompt(config.PasswordLabel); err != nil {
		return rec, err
	}

	var tagsString string

	if tagsString, err = line.Prompt(config.TagsLabel); err != nil {
		return rec, err
	}

	rec.Tags = tagsStringToArray(tagsString)

	if rec.Url, err = line.Prompt(config.URLLabel); err != nil {
		return rec, err
	}

	// make sure to close the line. Otherwise the following
	// scanner won't work.
	line.Close()

	log.Info("\nPlease insert Notes, end with EOF:")

	scanner := bufio.NewScanner(os.Stdin)

	text := ""

	for scanner.Scan() {
		input := scanner.Text()
		text += input + "\n"
	}

	err = scanner.Err()

	if err != nil {
		return rec, err
	}

	rec.Notes = text

	return rec, nil
}

// Print prints the content of the whole Masterfile to the console with DEBUG-level.
func (masterfile *MasterFile) Print(spacing int) {
	log.Debug("%*sMagic      : %s", spacing, "", masterfile.Magic)
	log.Debug("%*sMD5        : %s", spacing, "", masterfile.Md5)
	log.Debug("%*sGeneration : %d", spacing, "", masterfile.Generation)
}
