package services

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"wswagner.visualstudio.com/iotgate-server/v1/model"
)

func getPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyFile := os.Getenv("SIGNATURE_PRIVATE_KEYFILE")
	if privateKeyFile == "" {
		return nil, errors.New("no private signature key file provided")
	}
	_, err := os.Stat(privateKeyFile)
	if os.IsNotExist(err) {
		return nil, errors.New("specified private signature key file not found")
	}
	content, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("require rsa type key")
	}

	return rsaKey, nil
}

func GetPublicKeyPEM() ([]byte, error) {
	publicKeyFile := os.Getenv("SIGNATURE_PUBLIC_KEYFILE")
	if publicKeyFile == "" {
		return nil, errors.New("no public signature key file provided")
	}
	_, err := os.Stat(publicKeyFile)
	if os.IsNotExist(err) {
		return nil, errors.New("specified public signature key file not found")
	}
	return os.ReadFile(publicKeyFile)
}

func GetPublicKey() (*rsa.PublicKey, error) {
	content, err := GetPublicKeyPEM()
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("require rsa type key")
	}
	return rsaKey, nil
}

func TestPrivateKey() error {
	_, err := getPrivateKey()
	return err
}

func TestPublicKey() error {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return err
	}
	sig, _ := SignSHA256Hash(buf)
	err = VerifySHA256Signature(buf, sig)
	return err
}

func SignSHA256Hash(sha256hash []byte) ([]byte, error) {
	key, err := getPrivateKey()
	if err != nil {
		return nil, err
	}
	return rsa.SignPKCS1v15(nil, key, crypto.SHA256, sha256hash)
}

func VerifySHA256Signature(sha256hash []byte, sha256signature []byte) error {
	key, err := GetPublicKey()
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(key, crypto.SHA256, sha256hash, sha256signature)
}

func ComputeFirmwareSignature(meta model.FirmwareMetadata) ([]byte, error) {
	sha256Hash, err := model.DecodeMetaBytes(meta.SHAHash)
	if err != nil {
		return nil, err
	}

	signatureBytes, err := SignSHA256Hash(sha256Hash)
	if err != nil {
		return nil, err
	}

	// check verification via provided public key
	err = VerifySHA256Signature(sha256Hash, signatureBytes)
	if err != nil {
		return nil, err
	}

	return signatureBytes, nil
}
