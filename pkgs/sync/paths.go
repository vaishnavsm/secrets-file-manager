package sync

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/vaishnavsm/secrets-file-manager/pkgs/config"
)

func getUniqueValues[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}

	for _, val := range slice {
		if !seen[val] {
			seen[val] = true
			result = append(result, val)
		}
	}
	return result
}

type Path struct {
	CryptFile string
	PlainFile string
}

func GetPaths(config *config.Config) []Path {
	allPaths := []Path{}

	for _, path := range config.Paths {
		slog.Debug("syncing secrets file path", "path", path)

		files, err := filepath.Glob(path)
		if err != nil {
			slog.Error("error globbing path", "error", err)
			continue
		}

		for _, file := range files {
			if strings.HasSuffix(file, config.CryptSuffix) {
				allPaths = append(allPaths, Path{
					CryptFile: file,
					PlainFile: strings.TrimSuffix(file, config.CryptSuffix),
				})
			} else {
				allPaths = append(allPaths, Path{
					CryptFile: file + config.CryptSuffix,
					PlainFile: file,
				})
			}
		}
		cryptfiles, err := filepath.Glob(path + config.CryptSuffix)
		if err != nil {
			slog.Error("error globbing path", "error", err)
			continue
		}

		for _, file := range cryptfiles {
			allPaths = append(allPaths, Path{
				CryptFile: file,
				PlainFile: strings.TrimSuffix(file, config.CryptSuffix),
			})
		}
	}

	allPaths = getUniqueValues(allPaths)
	slog.Debug("got list of paths", "paths", allPaths)

	return allPaths
}
