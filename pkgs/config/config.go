package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Paths               []string `yaml:"paths"`
	CryptSuffix         string   `yaml:"crypt_suffix"`
	EnableReadFromCrypt bool     `yaml:"enable_read_from_crypt"`
	SecretMethod        string   `yaml:"secret_method"`
	PasswordFile        string   `yaml:"password_file,omitempty"`
	ForceSync           string   `yaml:"_force_sync,omitempty"`
	LogLevel            string   `yaml:"_log_level,omitempty"`
}

type Options struct {
	ForceSync string
	LogLevel  string
}

func CreateConfigFile() string {
	config := Config{
		Paths:               []string{"**/*.secrets.env", "*.secrets.env"},
		CryptSuffix:         ".secrets",
		SecretMethod:        "passwordfile",
		PasswordFile:        ".password",
		EnableReadFromCrypt: false,
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		slog.Error("error marshalling default config", "error", err)
		return ""
	}

	return string(data)
}

func LoadConfig(path string, options *Options) (*Config, error) {
	var config Config = Config{
		CryptSuffix:  ".secrets",
		SecretMethod: "passwordfile",
		PasswordFile: ".password",
		ForceSync:    options.ForceSync,
		LogLevel:     options.LogLevel,
	}

	slog.Debug("loading config", "path", path)

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	slog.Debug("unmarshalling config", "data", string(data))
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
