package utils

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"io/ioutil"
	"loki/config"
	"loki/log"
	pb "loki/storage"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// PromptPassword prompts the user for a password, optional twice and verifies equality if needed.
func PromptPassword(twice bool) ([]byte, error) {
	log.Info("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))

	if twice {
		log.Info("\nRe-enter Password again: ")
		bytePasswordAgain, err := terminal.ReadPassword(int(syscall.Stdin))

		if err != nil {
			return []byte{}, errors.New("problems reading second password from terminal")
		}

		if !bytes.Equal(bytePasswordAgain, bytePassword) {
			return []byte{}, errors.New("passwords do not match")
		}
	}

	if err != nil {
		return []byte{}, errors.New("problems reading first password from terminal")
	}

	log.Info("\n")
	return bytePassword, nil
}

// Hexdump provides a string with the hex-representation of the byte-array given in data.
func Hexdump(data []byte) string {
	return hex.EncodeToString(data)
}

// NormalizePath appends the .loki suffix on the path if its not allready present.
func NormalizePath(path string) string {
	if !strings.HasSuffix(path, config.FileSuffix) {
		return path + config.FileSuffix
	}
	return path
}

// VerifyDirectory makes sure the path given in dirname is present and a directory.
func VerifyDirectory(dirname string) bool {
	log.Debug("Verify directory: %s", dirname)
	fi, err := os.Stat(dirname)

	if err != nil {
		log.Error("Base directory does not exist : " + dirname)
		return false
	}

	if !fi.IsDir() {
		log.Error("Path given is not a directory: " + dirname)
		return false
	}
	log.Debug("Directory exists.")
	return true
}

// VerifyFile makes sure the path given in filename is present and is a file.
func VerifyFile(filename string) bool {
	fi, err := os.Stat(filename)

	if err != nil {
		log.Error("File does not exist : " + filename)
		return false
	}

	if fi.IsDir() {
		log.Error("Path given is a directory: " + filename)
		return false
	}
	return true
}

// CheckBase verifies that the system directory ( usually ~/.loki ) is present
// and accesible
func CheckBase(cfg config.Configuration) bool {
	base := cfg.SystemDirectory()

	if !VerifyDirectory(base) {
		log.Error("You might want to initialize the directoy using :")
		log.Error(" loki init")
		return false
	}
	return true
}

// InitBasedir creates a new system directory if needed and populates it with
// .master and .config file.
func InitBasedir(cfg config.Configuration) error {

	dirname := cfg.SystemDirectory()

	_, err := os.Stat(dirname)

	if err == nil {
		return errors.New("Base directory already exist : " + dirname)
	}

	err = os.Mkdir(dirname, os.ModePerm)

	if err != nil {
		return err
	}

	err = WriteNewMasterfile(cfg.GetMasterfilename())

	if err != nil {
		return err
	}

	err = createConfigfile(dirname, cfg.Gitmode)

	if err != nil {
		return err
	}

	if cfg.Gitmode {
		gitCmd := []string{"init", dirname}
		log.Debug("Running git command: %v", gitCmd)
		err = GitCommand(gitCmd)

		if err != nil {
			return err
		}

		err = os.Chdir(dirname)

		if err != nil {
			return err
		}

		if gitAddFile(config.MasterFilename); err != nil {
			return err
		}

		if gitAddFile(config.ConfigFilename); err != nil {
			return err
		}

		gitCommit := []string{"commit", "-a", "-m", "Initialized directory"}
		log.Debug("Running git command: %v", gitCommit)
		GitCommand(gitCommit)

		if err != nil {
			return err
		}
	}

	return nil
}

func gitAddFile(filename string) error {
	gitCmd := []string{"add", filename}
	log.Debug("Running git command: %v", gitCmd)
	return GitCommand(gitCmd)
}

func createConfigfile(dirname string, withGit bool) error {

	dst := dirname + string(os.PathSeparator) + config.ConfigFilename

	log.Debug("Create configfile: %s", dst)

	data := []byte(createConfigfileTemplate(withGit))

	return WriteFile(dst, data)
}

func createConfigfileTemplate(withGit bool) string {
	return fmt.Sprintf("[basic]\nLoglevel = INFO\nExternalEditor = false\nGitmode = %t\n", withGit)
}

// CopyFile copies a file content to a new location.
// From: https://stackoverflow.com/questions/21060945/simple-way-to-copy-a-file-in-golang
// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// GetBinaryPath returns the full qualified path of the binary actually running.
func GetBinaryPath() string {
	exe, _ := os.Executable()

	log.Debug("Executable: %s", exe)

	dir, _ := filepath.Abs(filepath.Dir(exe))
	return dir
}

// WriteFile stores the bytes provided in the byte-array data at the location given with path.
func WriteFile(path string, data []byte) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	n1, err := f.Write(data)

	if err != nil {
		return err
	}

	if n1 != len(data) {
		return errors.New("Could not write all data")
	}

	return nil
}

// Highlight searches for the searchstring in the string provided with text and
// highlights the matched string with terminal ASCII-codes for RED.
func Highlight(text string, searchstring string) string {
	info := color.New(color.FgRed).SprintFunc()

	idx := strings.Index(strings.ToLower(text), strings.ToLower(searchstring))

	if idx == -1 {
		return text
	}

	// return "|" + text[:idx] + "|" + info(searchstring) + "|" + text[idx+len(searchstring):] + "|"
	return text[:idx] + info(searchstring) + text[idx+len(searchstring):]
}

// Display displays the given lokifile-record at default column 0. Hides password if blind == true
func Display(rec *pb.Record, blind bool) {
	PrefixedDisplay(rec, 0, blind)
}

// PrefixedDisplay displays a lokifile record at the column provided by the spacing parameter.
// Hides password if blind == true
func PrefixedDisplay(rec *pb.Record, spacing int, blind bool) {
	PrefixedDisplayWithHighlighting(rec, spacing, "", blind)
}

// PrefixedDisplayWithHighlighting displays the given record at the coulunn given with spacing and
// highlights the searchstring if found. Hides password if blind == true
func PrefixedDisplayWithHighlighting(rec *pb.Record, spacing int, searchstring string, blind bool) {

	p := func(text string, searchstring string) string {
		if len(searchstring) > 0 {
			return Highlight(text, searchstring)
		}
		return text
	}

	log.Debug("%*s%s%s", spacing, "", config.MagicLabel, p(rec.Magic, searchstring))
	log.Debug("%*s%s%s", spacing, "", config.MD5Label, p(rec.Md5, searchstring))
	log.Info("%*s%s%s", spacing, "", config.TitleLabel, p(rec.Title, searchstring))
	log.Info("%*s%s%s", spacing, "", config.AccountLabel, p(rec.Account, searchstring))

	if blind {
		bLen := len(rec.Password)
		log.Info("%*s%s%s", spacing, "", config.PasswordLabel, strings.Repeat("*", bLen))
	} else {
		log.Info("%*s%s%s", spacing, "", config.PasswordLabel, p(rec.Password, searchstring))
	}

	if len(rec.Url) > 0 {
		log.Info("%*s%s%s", spacing, "", config.URLLabel, p(rec.Url, searchstring))
	}

	if len(rec.Notes) > 0 {
		log.Info("")

		longestLine := LongestLine(rec.Notes)
		separator := strings.Repeat("-", longestLine)

		log.Info("%*s%s", spacing, "", separator)
		log.Info("%*s%s", spacing, "", p(rec.Notes, searchstring))
		log.Info(separator)
	}
}

// CreateLeadingDirectories create the leading directories if they are missing
// to reach the file provided in filename.
func CreateLeadingDirectories(filename string) string {
	// check the filepath and create directories if needed
	dir, file := filepath.Split(filename)

	if len(dir) == 0 {
		log.Debug("Only filename given, no directory portion included")
	} else {

		dir := strings.TrimSuffix(dir, string(os.PathSeparator))
		log.Debug("Path given : " + dir)

		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Error("Error making directories: %v", err)
			return ""
		}

	}

	return file
}

// StartEditorWithData starts the editor located with the systems EDITOR environment variable
// with the contents provided in the data parameter. A temporary file with the content is
// created to achieve this.
func StartEditorWithData(data string) (string, error) {

	systemEditor := os.Getenv(config.LokiEditorEnv)

	if len(systemEditor) < 1 {
		return "", errors.New("Could not find editor environment variable: " + config.LokiEditorEnv)
	}
	log.Debug("Using editor from environment variable (%s) : %s", config.LokiEditorEnv, systemEditor)

	filename, err := StoreInTempfile(data)

	if err != nil {
		return "", err
	}

	log.Debug("Tempfilename for editing: %s", filename)

	defer os.Remove(filename)

	// if EDITOR is an alias, we gotta split it
	var params []string
	executable := ""

	if strings.Contains(systemEditor, " ") {
		parts := strings.Fields(systemEditor)
		executable = parts[0]
		params = parts[1:]
	} else {
		executable = systemEditor
	}

	params = append(params, filename)

	log.Debug("Executable : %s", executable)
	log.Debug("Params     : %v", params)

	cmd := exec.Command(executable, params...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return "", err
	}

	editedContent, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", err
	}

	return string(editedContent), nil
}

// StoreInTempfile stores the data provided in the parameter in a temporary file
// and returns the filename. Caller needs to cleanup any mess.
func StoreInTempfile(data string) (string, error) {

	f, err := ioutil.TempFile(os.TempDir(), "lokiEditBuffer")

	defer f.Close()

	if err != nil {
		return "", err
	}

	n1, err := f.Write([]byte(data))

	if err != nil {
		return "", err
	}

	if n1 != len(data) {
		return "", errors.New("Could not write all data")
	}

	name := f.Name()
	f.Close()

	return name, nil
}

// ExitSystemFailure exits the software with an predefined error-code
func ExitSystemFailure() {
	ExitSystemWithCode(config.ExitCodeFailure)
}

// ExitSystem with an go error type
func ExitSystem(err error) {
	if err != nil {
		log.Error("Exit system: %v", err)
		ExitSystemWithCode(config.ExitCodeFailure)
	}

	ExitSystemWithCode(config.ExitCodeOK)
}

// ExitSystemWithCode is the systems single exit point to the os.
// The given code will be the return code.
func ExitSystemWithCode(code int) {
	os.Exit(code)
}

// GitCommand runs the git subcommand given in the params array, connects stdin, -out, -err and waits for the commad to finish.
func GitCommand(params []string) error {
	cmd := exec.Command("git", params...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// LongestLine is supposed to get an multi-line string (\n terminated) and returns
// the width of the longest line in int.
func LongestLine(input string) int {
	lines := strings.Split(input, "\n")
	longest := 0

	for _, line := range lines {
		l := len(line)
		if l > longest {
			longest = l
		}
	}

	return longest
}

// CreateTempdirWithPrefx create a temporary directory in the systems
// default tempdir prefixed with the given prefix.
func CreateTempdirWithPrefx(prefix string) (string, error) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), prefix)

	if err != nil {
		return "", err
	}
	return tmpDir, nil
}
