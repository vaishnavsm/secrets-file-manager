package encrypt

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"os"

	"github.com/vaishnavsm/secrets-file-manager/pkgs/config"
)

func VersionToBytes(version uint16) []byte {
	return []byte{byte(version >> 8), byte(version)}
}
func VersionFromBytes(versionBytes []byte) (version uint16) {
	return uint16(versionBytes[0])<<8 | uint16(versionBytes[1])
}

func Encrypt(config *config.Config, path string) ([]byte, error) {
	plainData, err := os.ReadFile(path)
	if err != nil {
		slog.Error("error reading file", "error", err)
		return nil, err
	}

	method := config.SecretMethod

	var encryptedData []byte
	switch method {
	case "passwordfile":
		{
			encryptedData, err = encryptWithPasswordFile(config, plainData)
		}
	default:
		return nil, errors.New("unknown encryption method")
	}

	if err != nil {
		slog.Error("error encrypting data", "error", err)
		return nil, err
	}

	encodedData := base64.StdEncoding.EncodeToString(encryptedData)

	return []byte(encodedData), nil
}

func Decrypt(config *config.Config, path string) (plainData []byte, err error) {

	encodedEncryptedData, err := os.ReadFile(path)
	if err != nil {
		slog.Error("error reading file", "error", err)
		return nil, err
	}

	encryptedData, err := base64.StdEncoding.DecodeString(string(encodedEncryptedData))
	if err != nil {
		slog.Error("error decoding data", "error", err)
		return nil, err
	}

	method := config.SecretMethod

	switch method {
	case "passwordfile":
		return decryptWithPasswordFile(config, encryptedData)
	}

	return nil, errors.New("unknown encryption method")
}
