package services

import (
	"wswagner.visualstudio.com/iotgate-server/v1/model"
)

// DeviceDataService provides functions for retreiving and saving model.DeviceData entities
type DeviceDataService interface {
	GetDeviceData(deviceID string) (model.DeviceData, error)
	SetDeviceData(data model.DeviceData) error
}
