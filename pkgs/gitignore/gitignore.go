package gitignore

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/vaishnavsm/secrets-file-manager/pkgs/config"
	"github.com/vaishnavsm/secrets-file-manager/pkgs/sync"
)

func GenerateGitignore(config *config.Config) (err error) {

	gitignore, err := os.ReadFile(".gitignore")
	gitignoreLines := []string{}
	if err != nil {
		slog.Info("no .gitignore file found, creating one")
	} else {
		gitignoreLines = strings.Split(string(gitignore), "\n")
	}

	paths := sync.GetPaths(config)

	secretPaths := []string{
		config.PasswordFile,
	}

	for _, path := range paths {
		secretPaths = append(secretPaths, path.PlainFile)
	}

	if len(gitignoreLines) == 0 || gitignoreLines[len(gitignoreLines)-1] != "" {
		// Add a newline to the end of the file to prevent appending to the last line
		fmt.Println()
	}

	for _, path := range secretPaths {
		if !slices.Contains(gitignoreLines, path) {
			fmt.Println(path)
		}
	}

	return nil
}
