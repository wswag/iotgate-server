package services

import (
	"io"

	"github.com/wswag/iotgate-server/v1/model"
)

// FirmwareService provides methods to get metadata and download/upload firmware binary data
type FirmwareService interface {
	GetFirmwareMetadata(deviceID string) (model.FirmwareMetadata, error)
	SetFirmwareMetadata(data model.FirmwareMetadata) error
	UploadFirmware(deviceID string, src io.Reader, meta *model.FirmwareMetadata) error
	DownloadFirmware(deviceID string, dst io.Writer) error
	DownloadFirmwareChunk(deviceID string, start uint32, length uint32) ([]byte, error)
}
