package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/vaishnavsm/secrets-file-manager/pkgs/config"
)

const (
	PasswordFileEncryptionVersion uint16 = 1
)

func encryptWithPasswordFile(config *config.Config, plainData []byte) (encryptedData []byte, err error) {

	passwordFile := config.PasswordFile

	passwordFileStats, err := os.Stat(passwordFile)
	if err != nil {
		slog.Error("error getting password file stats", "error", err)
		return nil, err
	}
	if passwordFileStats.Mode()&0077 != 0 {
		slog.Error("password file is not secure. please make sure it is not readable by other users. run `chmod 600 " + passwordFile + "` or similar.")
		return nil, errors.New("password file is not secure")
	}

	passwordBytes, err := os.ReadFile(passwordFile)
	if err != nil {
		slog.Error("error reading password file", "error", err)
		return nil, err
	}

	password := strings.TrimSpace(string(passwordBytes))

	salt := make([]byte, 32)
	rand.Read(salt)

	key, err := pbkdf2.Key(sha256.New, password, salt, 4096, 32)
	if err != nil {
		slog.Error("error generating key with pbkdf2", "error", err)
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		slog.Error("error creating aes cipher", "error", err)
		return nil, err
	}

	stream, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("error creating gcm", "error", err)
		return nil, err
	}

	nonce := make([]byte, stream.NonceSize())
	rand.Read(nonce)

	encryptedData = stream.Seal(nil, nonce, plainData, nil)

	return append(append(append(VersionToBytes(PasswordFileEncryptionVersion), salt...), nonce...), encryptedData...), nil
}

func decryptWithPasswordFile(config *config.Config, encryptedData []byte) (plainData []byte, err error) {

	passwordFile := config.PasswordFile

	passwordBytes, err := os.ReadFile(passwordFile)
	if err != nil {
		slog.Error("error reading password file", "error", err)
		return nil, err
	}

	password := strings.TrimSpace(string(passwordBytes))

	version := VersionFromBytes(encryptedData[:2])

	if version != PasswordFileEncryptionVersion {
		slog.Error("unsupported version", "version", version, "expected", PasswordFileEncryptionVersion)
		return nil, errors.New("unsupported version " + strconv.Itoa(int(version)))
	}

	salt := encryptedData[2:34]

	key, err := pbkdf2.Key(sha256.New, password, salt, 4096, 32)
	if err != nil {
		slog.Error("error generating key with pbkdf2", "error", err)
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		slog.Error("error creating aes cipher", "error", err)
		return nil, err
	}

	stream, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("error creating gcm", "error", err)
		return nil, err
	}

	nonceEnd := 34 + stream.NonceSize()
	nonce := encryptedData[34:nonceEnd]
	encryptedData = encryptedData[nonceEnd:]

	plainData, err = stream.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		slog.Error("error decrypting data", "error", err)
		return nil, err
	}

	return plainData, nil

}
