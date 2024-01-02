package services

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"

	"github.com/wswag/iotgate-server/v1/model"
)

type ESP32MetadataExtractor struct {
}

type ESP32ImageHeader struct {
	Magic              uint8
	Segment_count      uint8
	Spi_mode           uint8
	Spi_speed_and_size uint8
	Entry_addr         uint32
	Wp_pin             uint8
	Spi_pin_drv        [3]uint8
	Chip_id            uint16
	Min_chip_rev       uint8
	Reserved           [8]uint8
	Hash_appended      uint8
}

type ESP32ImageSegmentHeader struct {
	Load_addr uint32
	Data_len  uint32
}

type ESP32AppDescr struct {
	Magic_word     uint32
	Secure_version uint32
	Reserv1        [2]uint32
	Version        [32]byte
	Project_name   [32]byte
	Time           [16]byte
	Date           [16]byte
	Idf_ver        [32]byte
	App_elf_sha256 [32]byte
	Reserv2        [20]uint32
}

func (e ESP32MetadataExtractor) ExtractMeta(firmware []byte, meta *model.FirmwareMetadata) {
	r := bytes.NewReader(firmware)
	imgHeader := ESP32ImageHeader{}
	imgSegmHeader := ESP32ImageSegmentHeader{}
	esph := ESP32AppDescr{}

	binary.Read(r, binary.LittleEndian, &imgHeader)
	binary.Read(r, binary.LittleEndian, &imgSegmHeader)
	binary.Read(r, binary.LittleEndian, &esph)

	log.Println("Hash appended: ", imgHeader.Hash_appended)
	if imgHeader.Hash_appended != 0 {
		// replace computed hash with appended hash
		log.Println("Extracting new Hash")
		appendedHash := (firmware[len(firmware)-32:])
		computedHash := sha256.Sum256(firmware[:len(firmware)-32])
		if !bytes.Equal(appendedHash, computedHash[:]) {
			log.Println("warning: appended hash differs from computed hash!")
		}
		meta.SHAHash = model.EncodeMetaBytes(computedHash[:])
	} else {
		computedHash := sha256.Sum256(firmware)
		meta.SHAHash = model.EncodeMetaBytes(computedHash[:])
	}

	log.Println("Extracted metadata:")
	log.Println(esph.Magic_word)
	log.Println(esph.Secure_version)
	log.Println(string(esph.Version[:]))
	log.Println(string(esph.Project_name[:]))
	log.Println(string(esph.Time[:]))
	log.Println(string(esph.Date[:]))
	log.Println(string(esph.Idf_ver[:]))
	log.Println(string(esph.App_elf_sha256[:]))
}
