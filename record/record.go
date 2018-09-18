package record

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"loki/config"
	"loki/crypto"
	"loki/log"
	pb "loki/storage"
	"loki/utils"
)

// This is the basic binary header of every loki-file
//
// Magic      : 4c 4f 4b 49     :  4 : "LOKI" Magic Header
// Version    : 00 00 00 01     :  4 : v1 - Protocol/Format version
// Generation : 00 00 00 17     :  4 : Generation number of Masterpassword
// Size       : 00 00 00 00     :  4 : Size of encrypted payload
// MD5 Hash   : 16 Bytes        : 16 : md5sum of encrypted payload
//
// Data       : .........       : Variable-sized encrypted payload
const (
	LokiFormatVersion uint32 = 1
	LokiHeaderSize    int    = 32

	MagicValue1 byte = 0x4c
	MagicValue2 byte = 0x4f
	MagicValue3 byte = 0x4b
	MagicValue4 byte = 0x49
)

// DataFileHeader is the, mmmh, well the header of any lokifile
type DataFileHeader struct {
	FormatVersion uint32
	Generation    uint32
	PayloadSize   uint32
	PayloadMD5    []byte
}

// Print prints the DataFileHeader prefixed with a number of spaces provided by the column parameter.
func (hdr *DataFileHeader) Print(column int) {
	log.Debug("%*sFormat version : %d", column, "", hdr.FormatVersion)
	log.Debug("%*sGeneration     : %d", column, "", hdr.Generation)
	log.Debug("%*sPayload Size   : %d", column, "", hdr.PayloadSize)
	log.Debug("%*sPayload MD5    : %s", column, "", utils.Hexdump(hdr.PayloadMD5))
}

var engine = crypto.NewEngine()

// ComputeInnerMd5 returns a hex-encoded string of the md5 hash of all fields for the provided record.
func ComputeInnerMd5(rec pb.Record) string {
	return utils.Hexdump(crypto.GetStringMD5(rec.Title + rec.Account + rec.Password + rec.Url + rec.Notes))
}

// WriteRecord saves the given record using the given key to the location given with path.
func WriteRecord(path string, generation uint32, key []byte, rec pb.Record) error {
	// Adding Magic and inner-payload md5
	rec.Magic = config.InnerMagic
	rec.Md5 = ComputeInnerMd5(rec)

	serialized, err := proto.Marshal(&rec)

	if err != nil {
		return err
	}

	encryptedPayload, err := engine.Encrypt(serialized, key)

	if err != nil {
		return err
	}

	hdr := createHeaderForPayload(encryptedPayload, generation)

	return utils.WriteFile(path, append(hdr[:], encryptedPayload[:]...))
}

func createHeaderForPayload(payload []byte, generation uint32) []byte {
	header := make([]byte, LokiHeaderSize-16)

	header[0] = MagicValue1
	header[1] = MagicValue2
	header[2] = MagicValue3
	header[3] = MagicValue4

	binary.BigEndian.PutUint32(header[4:], LokiFormatVersion)
	binary.BigEndian.PutUint32(header[8:], generation)
	binary.BigEndian.PutUint32(header[12:], uint32(len(payload)))

	checksum := crypto.ComputeMD5checksum(payload)
	return append(header[:], checksum...)
}

// LoadRecord returns a valid record if it could decrypt with the given key the file provided with filename.
func LoadRecord(filename string, key []byte) (*pb.Record, *DataFileHeader, error) {

	f, err := os.Open(filename)

	if err != nil {
		return &pb.Record{}, &DataFileHeader{}, errors.New("file not found")
	}

	defer f.Close()

	header := make([]byte, LokiHeaderSize)

	n1, err := f.Read(header)

	if err != nil {
		return &pb.Record{}, &DataFileHeader{}, fmt.Errorf("error reading header: %v", err)
	}

	if n1 != LokiHeaderSize {
		return &pb.Record{}, &DataFileHeader{}, fmt.Errorf("header corrupted")
	}

	hdr, err := parseHeader(header)

	if err != nil {
		return &pb.Record{}, &DataFileHeader{}, fmt.Errorf("error parsing header: %v", err)
	}

	// hdr.Print()

	payload, err := readPayload(f, hdr.PayloadSize)

	if err != nil {
		return &pb.Record{}, &DataFileHeader{}, err
	}

	// decrypt
	decryptedPayload, err := engine.Decrypt(payload, key)

	rec := &pb.Record{}

	if err != nil {
		return rec, &DataFileHeader{}, errors.New("unable to decrypt, password?")
	}

	err = proto.Unmarshal(decryptedPayload, rec)

	if err != nil {
		return rec, &DataFileHeader{}, errors.New("error unmarshaling payload")
	}

	if !crypto.VerifyMD5(payload, hdr.PayloadMD5) {
		return rec, &DataFileHeader{}, errors.New("md5 checksum incorrect")
	}

	if rec.Md5 != ComputeInnerMd5(*rec) {
		return rec, &DataFileHeader{}, errors.New("inner MD5 checksum incorrect")
	}

	if rec.Magic != config.InnerMagic {
		return rec, &DataFileHeader{}, errors.New("inner Magic not correct")
	}

	return rec, &hdr, nil
}

func readPayload(f *os.File, payloadsize uint32) ([]byte, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, errors.New("Could not stat file")
	}

	filesize := fi.Size()

	// verify given payload size against headersize + file-length

	if int64(LokiHeaderSize)+int64(payloadsize) != filesize {
		return nil, fmt.Errorf("Sizes do not match. Headersize: %d, Given payload: %d, Filesize: %d",
			LokiHeaderSize, payloadsize, filesize)
	}

	buffer := make([]byte, payloadsize)
	bytesRead, err := f.Read(buffer)

	if bytesRead != int(payloadsize) {
		return nil, errors.New("Could not read Payload")
	}
	return buffer, nil
}

func verifyMagic(magic []byte) bool {
	if magic[0] != MagicValue1 ||
		magic[1] != MagicValue2 ||
		magic[2] != MagicValue3 ||
		magic[3] != MagicValue4 {
		return false
	}
	return true
}

func parseHeader(header []byte) (DataFileHeader, error) {

	if len(header) != LokiHeaderSize {
		return DataFileHeader{}, errors.New("header size incorrect")
	}

	if !verifyMagic(header) {
		return DataFileHeader{}, errors.New("header magic value incorrect, must be " + config.InnerMagic)
	}

	hdr := DataFileHeader{}

	hdr.FormatVersion = binary.BigEndian.Uint32(header[4:8])
	hdr.Generation = binary.BigEndian.Uint32(header[8:12])
	hdr.PayloadSize = binary.BigEndian.Uint32(header[12:16])

	hdr.PayloadMD5 = make([]byte, 16)
	copy(hdr.PayloadMD5, header[16:])

	return hdr, nil
}
