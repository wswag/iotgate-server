package services

import (
	"crypto/sha256"
	"errors"
	"io"

	"wswagner.visualstudio.com/iotgate-server/v1/model"
)

// MFFirmwareService implements FirmwareService based on a ManagedFileService
type MFFirmwareService struct {
	FileService ManagedDocumentService
}

const firmwareTopic = "Firmware"
const firmwareMetaTopic = "Firmware.Meta"

// GetFirmwareMetadata returns the persisted metadata related to deviceID
func (m MFFirmwareService) GetFirmwareMetadata(deviceID string) (model.FirmwareMetadata, error) {
	result := model.FirmwareMetadata{}
	objSrv := ObjectDocumentService{DocSrv: m.FileService}
	err := objSrv.OpenObject(firmwareMetaTopic, deviceID, &result)
	return result, err
}

// SetFirmwareMetadata allows to override metadata
func (m MFFirmwareService) SetFirmwareMetadata(data model.FirmwareMetadata) error {
	objSrv := ObjectDocumentService{DocSrv: m.FileService}
	err := objSrv.CreateObject(firmwareMetaTopic, data.DeviceID, data)
	return err
}

// UploadFirmware persists the content in src as the new firmware for deviceID
func (m MFFirmwareService) UploadFirmware(deviceID string, src io.Reader, meta *model.FirmwareMetadata) error {
	f, err := m.FileService.Create(firmwareTopic, deviceID)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, src)
	if err != nil {
		return err
	}

	// TODO: dont load whole conent, instead use a stream
	fwFile, _ := m.FileService.Open(firmwareTopic, deviceID)
	bytes, _ := io.ReadAll(fwFile)
	meta.SHAHash = ""
	e := ESP32MetadataExtractor{}
	e.ExtractMeta(bytes, meta)

	if meta.SHAHash == "" {
		// not extracted from metadata so compute the hash
		hash := sha256.Sum256(bytes)
		meta.SHAHash = model.EncodeMetaBytes(hash[:])
	}

	return err
}

// DownloadFirmware copies the related firmware to the given io.Writer dst
func (m MFFirmwareService) DownloadFirmware(deviceID string, dst io.Writer) error {
	f, err := m.FileService.Open(firmwareTopic, deviceID)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(dst, f)
	return err
}

// DownloadFirmwareChunk returns the binary chunk specified by start and length
func (m MFFirmwareService) DownloadFirmwareChunk(deviceID string, start uint32, length uint32) ([]byte, error) {
	f, err := m.FileService.Open(firmwareTopic, deviceID)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	stat, _ := f.Stat()

	if length == 0 {
		length = uint32(stat.Size()) - start
	}

	buf := make([]byte, length)

	if int64(start+length) > stat.Size() {
		return buf, errors.New("invalid length")
	}

	_, err = f.Seek(int64(start), 0)
	if err != nil {
		return buf, err
	}

	_, err = f.Read(buf)
	return buf, err
}
