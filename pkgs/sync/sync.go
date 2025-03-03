package sync

import (
	"log/slog"
	"os"

	"github.com/vaishnavsm/secrets-file-manager/pkgs/config"
	"github.com/vaishnavsm/secrets-file-manager/pkgs/encrypt"
)

func Sync(config *config.Config) (err error) {

	slog.Debug("syncing secrets files")

	paths := GetPaths(config)

	for _, path := range paths {
		err = SyncFile(config, path)
		if err != nil {
			slog.Error("error syncing file", "error", err)
		}
	}

	return err
}

func SyncFile(config *config.Config, path Path) error {

	slog.Debug("syncing secrets file", "path", path)

	isCryptToFileSync := getIsCryptToFileSync(config, path.CryptFile, path.PlainFile)
	if isCryptToFileSync {
		slog.Debug("crypt file is newer than file, syncing crypt to file")
		err := syncCryptToFile(config, path.CryptFile, path.PlainFile)
		if err != nil {
			slog.Error("error syncing crypt to file", "error", err)
			return err
		}
	} else {
		slog.Debug("crypt file is not newer than file, syncing file to crypt")
		err := syncFileToCrypt(config, path.PlainFile, path.CryptFile)
		if err != nil {
			slog.Error("error syncing file to crypt", "error", err)
			return err
		}
	}

	return nil
}

func getIsCryptToFileSync(config *config.Config, cryptFile string, path string) bool {

	if config.ForceSync == "from-crypt" {
		return true
	}

	if config.ForceSync == "to-crypt" {
		return false
	}

	// does crypt exist?
	cryptExists, err := os.Stat(cryptFile)
	if err != nil {
		slog.Debug("error checking if crypt file exists", "error", err)
		return false
	}

	// does file exist?
	fileExists, err := os.Stat(path)
	if err != nil {
		slog.Debug("error checking if file exists", "error", err)
		return true
	}

	if !config.EnableReadFromCrypt {
		return false
	}

	// is crypt newer than file?
	if cryptExists.ModTime().After(fileExists.ModTime()) {
		return true
	}

	return false
}

func syncCryptToFile(config *config.Config, cryptFile string, path string) (err error) {

	plainData, err := encrypt.Decrypt(config, cryptFile)
	if err != nil {
		slog.Error("error decrypting file", "error", err)
		return err
	}

	err = os.WriteFile(path, plainData, 0600)
	if err != nil {
		slog.Error("error writing plain data to file", "error", err)
		return err
	}

	return nil
}

func syncFileToCrypt(config *config.Config, path string, cryptFile string) (err error) {

	encryptedData, err := encrypt.Encrypt(config, path)
	if err != nil {
		slog.Error("error encrypting file", "error", err)
		return err
	}

	err = os.WriteFile(cryptFile, encryptedData, 0600)
	if err != nil {
		slog.Error("error writing encrypted data to file", "error", err)
		return err
	}

	return nil
}
