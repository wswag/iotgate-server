package services

import (
	"github.com/wswag/iotgate-server/v1/model"
)

// DeviceDataService provides functions for retreiving and saving model.DeviceData entities
type DeviceDataService interface {
	GetDeviceData(deviceID string) (model.DeviceData, error)
	SetDeviceData(data model.DeviceData) error
}
