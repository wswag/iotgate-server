package services

import (
	"wswagner.visualstudio.com/iotgate-server/v1/model"
)

// MFDeviceDataService implements the DeviceDataService interface with a ManagedDocumentService
type MFDeviceDataService struct {
	FileService ManagedDocumentService
}

const deviceDataTopic = "DeviceData"

// GetDeviceData returns the model.DeviceData object related to deviceID
func (m MFDeviceDataService) GetDeviceData(deviceID string) (model.DeviceData, error) {
	result := model.DeviceData{}
	objSrv := ObjectDocumentService{DocSrv: m.FileService}
	err := objSrv.OpenObject(deviceDataTopic, deviceID, &result)
	return result, err
}

// SetDeviceData stores the given model.DeviceData
func (m MFDeviceDataService) SetDeviceData(data model.DeviceData) error {
	objSrv := ObjectDocumentService{DocSrv: m.FileService}
	err := objSrv.CreateObject(deviceDataTopic, data.DeviceID, data)
	return err
}
