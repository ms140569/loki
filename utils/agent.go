package utils

import (
	"errors"
	"loki/config"
	"loki/crypto"
	"loki/log"
	"net"
	"os"
	"os/exec"
)

// GetMasterkey tries to get the masterkey from the loki-agent running in the background.
func GetMasterkey(twice bool) ([]byte, error) {
	key, err := GetMasterkeyWithAgent(twice, true)

	if err != nil {
		return []byte{}, errors.New("Problem getting Masterkey")
	}

	return key, nil
}

// GetMasterkeyWithAgent tries to get the masterkey possibly from the agent or prompting the user once or twice
// according the twice flag.
func GetMasterkeyWithAgent(twice bool, withAgent bool) ([]byte, error) {

	if withAgent {
		if key, err := askAgent(); err == nil {
			log.Debug("Fine, got key from agent : " + Hexdump(key))
			return key, nil
		}
	}

	kdf := crypto.NewKeyDerivator()
	password, err := PromptPassword(twice)

	if err != nil {
		return []byte{}, errors.New("Problem prompting password")
	}

	key := kdf(password)
	return key, nil
}

func askAgent() ([]byte, error) {

	socketFile := config.GetSocketfilePath()

	if _, err := os.Stat(socketFile); err != nil {
		return []byte{}, errors.New("Socketfile not found")
	}

	c, err := net.Dial("unix", socketFile)

	if err != nil {
		return []byte{}, errors.New("Dial error")
	}

	defer c.Close()

	key := make([]byte, config.KeyLength)

	_, err = c.Write([]byte(config.RequestMagic))

	if err != nil {
		return []byte{}, errors.New("Error writing magic")
	}

	n, err := c.Read(key[:])

	if err != nil {
		return []byte{}, errors.New("Read error")
	}

	if n != config.KeyLength {
		return []byte{}, errors.New("Could not read all bytes from socket, but only : " + string(n))
	}

	return key, nil
}

// ShutdownAgent stops the background agent by sending it a stop command on the unix domain socket.
func ShutdownAgent() error {

	socketFile := config.GetSocketfilePath()

	if _, err := os.Stat(socketFile); err != nil {
		return err
	}

	c, err := net.Dial("unix", socketFile)

	if err != nil {
		return err
	}

	defer c.Close()

	_, err = c.Write([]byte(config.ShutdownMagic))

	if err != nil {
		return err
	}

	return nil
}

// SetupKeyAgent starts the background daemon to hold the systems key and passes the
// key on stdin  to the daemon.
func SetupKeyAgent(key []byte) error {
	return SetupKeyAgentWithBinpath(key, GetBinaryPath())
}

// SetupKeyAgentWithBinpath starts the background daemon to hold the systems key and passes the
// key on stdin  to the daemon. In addition one can provide the binpath. This is used for testing
// since the binarypath could not be derived from the main binary in this case.
func SetupKeyAgentWithBinpath(key []byte, binpath string) error {
	if _, err := os.Stat(config.GetSocketfilePath()); err == nil {
		log.Debug("Socketfile found, bail out.")
		return nil
	}

	cmd := exec.Command(binpath + string(os.PathSeparator) + "loki-agentd")
	childStdin, err := cmd.StdinPipe()

	if err != nil {
		return err
	}

	err = cmd.Start()

	if err != nil {
		return err
	}

	log.Debug("Sucessfully started agent, Path : " + cmd.Path)

	n, err := childStdin.Write(key)

	if err != nil {
		return err
	}

	if n != config.KeyLength {
		return errors.New("Could not write all bytes to pipe")
	}

	return nil
}
