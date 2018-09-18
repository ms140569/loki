package utils

import (
	"loki/config"
	"loki/crypto"
	pb "loki/storage"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"os"
)

// WriteNewMasterfile stores a brand new created masterfile at given path
func WriteNewMasterfile(path string) error {
	// there is no number zero, we always start with 1
	masterfile := createMasterfile(1)

	masterfile.Print(0)

	serialized, err := proto.Marshal(masterfile)

	if err != nil {
		return err
	}

	return WriteFile(path, serialized)

}

// RaiseGenerationInMasterfile loads masterfile given with path, increases the generation number by one and
// stores the file again.
func RaiseGenerationInMasterfile(path string) error {
	masterfile, err := LoadMasterfile(path)

	if err != nil {
		return errors.New("could not load masterfile")
	}

	masterfile = createMasterfile(masterfile.Generation + 1)

	serialized, err := proto.Marshal(masterfile)

	if err != nil {
		return err
	}

	return WriteFile(path, serialized)
}

func createMasterfile(generation uint32) *pb.MasterFile {

	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, generation)

	masterfile := pb.MasterFile{}
	masterfile.Magic = config.InnerMagic
	masterfile.Md5 = Hexdump(crypto.ComputeMD5checksum(b))
	masterfile.Generation = generation

	return &masterfile
}

// LoadMasterfile returns a valid systems masterfile located at path.
func LoadMasterfile(path string) (*pb.MasterFile, error) {

	f, err := os.Open(path)

	if err != nil {
		return &pb.MasterFile{}, errors.New("file not found")
	}

	defer f.Close()

	// Stat the path for the size:

	fi, err := os.Stat(path)

	if err != nil {
		return &pb.MasterFile{}, errors.New("could not Stat file")
	}

	masterfile := &pb.MasterFile{}

	buffer := make([]byte, int(fi.Size()))

	_, err = f.Read(buffer)

	if err != nil {
		return &pb.MasterFile{}, errors.New("could not read masterfile")
	}

	err = proto.Unmarshal(buffer, masterfile)

	if err != nil {
		return &pb.MasterFile{}, errors.New("could not unmarshal masterfile")
	}

	return masterfile, nil
}
