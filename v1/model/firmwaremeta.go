package model

// FirmwareMetadata contains basic information about an uploaded firmware for a specific
// device, including a signature
type FirmwareMetadata struct {
	DeviceID  string
	Timestamp int64
	Iteration uint16
	Size      uint32
	SHAHash   string
	Signature string
}
