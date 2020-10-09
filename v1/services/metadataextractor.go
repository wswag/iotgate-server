package services

import (
	"wswagner.visualstudio.com/iotgate-server/v1/model"
)

// MetadataExtractor allows implementing services that can extract metadata from firmware binary data
type MetadataExtractor interface {
	ExtractMeta(firmware []byte, meta *model.FirmwareMetadata)
}
