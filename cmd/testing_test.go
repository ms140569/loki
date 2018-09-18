package cmd

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"loki/config"
	"loki/log"
	"loki/subcommand"
	"loki/utils"
)

const TESTDATA = "../data/test/minimal"

var cfg config.Configuration
var cmd subcommand.Subcommand
var tmpDir string
var startDir string
var flagBundle config.FlagBundle

func getTestLoglevel() string {
	logLevel := os.Getenv(config.LokiLoglevelEnv)

	if len(logLevel) > 0 {
		return logLevel
	}
	return "Fatal"
}

const SEP = string(os.PathSeparator)

func TBASE() string {
	return tmpDir + string(os.PathSeparator)
}

func TestMain(m *testing.M) {

	ll := getTestLoglevel()
	fmt.Printf("Rung with loglevel: %s\n", ll)

	log.SetSystemLogLevelFromString(ll)

	flagBundle = config.ParseFlags()

	key, err := hex.DecodeString("f54d6aba8329dea96d4b3daa8caaa05e06bd10c246a40d510d2feb3e73b620bb")

	if err != nil {
		log.Error("Error decoding key")
		return
	}

	cwd, _ := os.Getwd()
	log.Debug("PWD: %s", cwd)

	startDir = cwd

	if err = utils.SetupKeyAgentWithBinpath(key, "../bin"); err != nil {
		log.Error("Problem setting-up keyagent for tests: %v", err)
		return
	}

	time.Sleep(1 * time.Second)

	code := m.Run()
	utils.ShutdownAgent()

	time.Sleep(1 * time.Second)
	os.Exit(code)

}

func SetupTest(t *testing.T) func() {
	setupTestBottom()
	return teardownTest
}

func setupTestBottom() {

	var err error

	tmpDir, err = ioutil.TempDir(os.TempDir(), "loki_testing_basedir")

	if err != nil {
		log.Error("Problem creating Testdir: %v", err)
		return
	}

	log.Debug("Createing Tempdir: %s", tmpDir)

	var info os.FileInfo

	info, err = os.Stat(startDir + "/" + TESTDATA)

	if err != nil {
		cwd, e2 := os.Getwd()

		if e2 != nil {
			log.Error("Error stating cwd: %v", e2)
		}

		log.Error("Error accessing template dir: %v, cwd: %s", err, cwd)
		return
	}

	utils.Copy(startDir+"/"+TESTDATA, tmpDir, info)
	setupReadOnly(tmpDir)
}

func teardownTest() {
	log.Debug("Deleteing tempdir: %s", tmpDir)

	err := os.RemoveAll(tmpDir)

	if err != nil {
		log.Error("Problems removing tempdir: %v", err)
	}
}

func setupReadOnly(base string) {
	cfg = config.New(flagBundle)
	cfg.Loglevel = getTestLoglevel()
	cfg.SetSystemDirectory(base)

	cmd = subcommand.Subcommand{}

	sysdir := cfg.SystemDirectory()
	_ = os.Chdir(sysdir)
	log.Debug("Setup test envrionment : %s", base)
}
