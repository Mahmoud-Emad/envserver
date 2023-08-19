package internal

import (
	"io"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Database DatabaseConfig `toml:"database"`
	Server   ServerConfig   `toml:"server"`
}

type ServerConfig struct {
	Host            string `toml:"host"`
	Port            int    `toml:"port"`
	JWTSecretKey    string `toml:"jwtSecretKey"`
	ShutdownTimeout int    `toml:"shutdownTimeout"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int64  `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
}

// Read the config file.
func ReadConfigFromFile(path string) (Config, error) {
	config := Config{}
	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		return Config{}, cantLoadConfigFileError
	}
	err = config.validateConfig()
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

// Read the config from a string.
func ReadConfigFromString(content string) (Config, error) {
	return ReadConfigFromReader(strings.NewReader(content))
}

// Read the config from a reader.
func ReadConfigFromReader(r io.Reader) (Config, error) {
	config := Config{}
	_, err := toml.NewDecoder(r).Decode(&config)
	if err != nil {
		return Config{}, cantDecodeConfigError
	}
	err = config.validateConfig()
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
