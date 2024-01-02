package services

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// ManagedDocumentFileService provides functions for reading and writing documents identified by a topic and a key
type ManagedDocumentFileService struct {
	Basepath       string
	StoragePattern string
}

func (m ManagedDocumentFileService) getFilename(topic string, key string) string {
	// use base64 encoding to prevent file system attacs by names
	topicPart := strings.ReplaceAll(base64.URLEncoding.EncodeToString([]byte(topic)), "=", "")
	keyPart := strings.ReplaceAll(base64.URLEncoding.EncodeToString([]byte(key)), "=", "")
	_, err := os.Stat(m.Basepath)
	if os.IsNotExist(err) {
		os.MkdirAll(m.Basepath, os.ModePerm)
	}
	return m.Basepath + fmt.Sprintf("/"+m.StoragePattern, topicPart, keyPart)
}

// Open will open a stream to the data residing under the given topic and key
func (m ManagedDocumentFileService) Open(topic string, key string) (*os.File, error) {
	fn := m.getFilename(topic, key)
	return os.Open(fn)
}

// Create will create a stream for writing data to
func (m ManagedDocumentFileService) Create(topic string, key string) (*os.File, error) {
	fn := m.getFilename(topic, key)
	return os.Create(fn)
}

// Exists checks whether the given topic/key combination exists
func (m ManagedDocumentFileService) Exists(topic string, key string) bool {
	fn := m.getFilename(topic, key)
	_, err := os.Stat(fn)
	return err == nil
}
