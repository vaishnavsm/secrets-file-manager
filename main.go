package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/vaishnavsm/secrets-file-manager/pkgs/config"
	"github.com/vaishnavsm/secrets-file-manager/pkgs/gitignore"
	"github.com/vaishnavsm/secrets-file-manager/pkgs/sync"
)

var (
	configFile string
	logLevel   string
	forceSync  string
)

func setLogLevel(level string) {
	loglevel := slog.LevelInfo
	switch level {
	case "debug":
		loglevel = slog.LevelDebug
	case "info":
		loglevel = slog.LevelInfo
	case "warn":
		loglevel = slog.LevelWarn
	case "error":
		loglevel = slog.LevelError
	}
	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: loglevel})
	slog.SetDefault(slog.New(h))
}

func init() {
	flag.StringVar(&configFile, "config", ".secrets-file-manager.yaml", "path to config file")
	flag.StringVar(&forceSync, "force-sync", "", "force sync in a particular direction. options: from-crypt, to-crypt")
	flag.StringVar(&logLevel, "log-level", "error", "log level, can be debug, info, warn, error")
}

func main() {

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	flag.Parse()

	if logLevel == "" {
		if os.Getenv("LOG_LEVEL") != "" {
			logLevel = os.Getenv("LOG_LEVEL")
		}
	}

	setLogLevel(logLevel)

	if os.Args[1] == "init" {
		handleInit()
		os.Exit(0)
	}

	config, err := config.LoadConfig(configFile, &config.Options{ForceSync: forceSync, LogLevel: logLevel})
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	slog.Debug("config loaded", "config", config)

	switch os.Args[1] {
	case "gen-gitignore":
		handleGenGitignore(config)
	case "sync":
		handleSync(config)
	case "watch":
		handleWatch(config)
	case "help":
		printUsage()
		os.Exit(1)
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	programParts := strings.Split(os.Args[0], "/")
	programName := programParts[len(programParts)-1]
	fmt.Println("Usage: " + programName + " [options] <command>")
	fmt.Println("\nCommands:")
	fmt.Println("  init  			Create config file template")
	fmt.Println("  gen-gitignore  	Generate .gitignore file")
	fmt.Println("  sync          	Sync secrets files")
	fmt.Println("  watch         	Watch for changes")
	fmt.Println("  help          	Show this help message")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
}

func handleInit() {
	fmt.Println("# pipe to .secrets-file-manager.yaml")
	fmt.Println(config.CreateConfigFile())
}
func handleGenGitignore(config *config.Config) {
	gitignore.GenerateGitignore(config)
}

func handleSync(config *config.Config) {
	slog.Debug("syncing secrets files")
	sync.Sync(config)
	fmt.Println("Done!")
}

func handleWatch(config *config.Config) {
	slog.Error("watch not implemented")
	os.Exit(1)
}
