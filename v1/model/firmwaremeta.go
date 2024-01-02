package model

import (
	"encoding/base64"
)

// FirmwareMetadata contains basic information about an uploaded firmware for a specific
// device, including a signature
type FirmwareMetadata struct {
	DeviceID  string // the device id / alias
	Timestamp int64  // the timestamp of upload
	Iteration uint16 // the firmware version on iotgate-server
	Size      uint32 // the firmware size
	SHAHash   string // the sha256 hash of the firmware
	Signature string // the RSA signature of the iotgate-server, verifiable via its public key
}

func EncodeMetaBytes(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}

func DecodeMetaBytes(str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(str)
}
