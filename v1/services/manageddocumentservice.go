package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ManagedDocumentService provides streaming access to persistent data stored by topic/key paradigm
type ManagedDocumentService interface {
	Open(topic string, key string) (*os.File, error) // TODO!
	Create(topic string, key string) (*os.File, error)
	Exists(topic string, key string) bool
}

// ObjectDocumentService provides convenience functions to store objects within a ManagedDocumentService
type ObjectDocumentService struct {
	DocSrv ManagedDocumentService
}

// OpenObject is a convenience function to directly unmarshal a sructure from a given key/value file.
// The structure is unmarshaled as JSON and may be inversely saved via CreateObject before
func (m *ObjectDocumentService) OpenObject(topic string, key string, dst interface{}) error {
	f, err := m.DocSrv.Open(topic, key)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &dst)
	return err
}

// CreateObject is a convenience function to directly marshal a sructure to a given key/value file.
// The structure is marshaled as JSON and may be inversely loaded via OpenObject afterwards
func (m *ObjectDocumentService) CreateObject(topic string, key string, obj interface{}) error {
	f, err := m.DocSrv.Create(topic, key)
	if err != nil {
		return err
	}
	defer f.Close()

	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	return err
}
