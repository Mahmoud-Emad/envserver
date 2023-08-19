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
	JWTSecretKey    string `toml:"jwt_secret_key"`
	ShutdownTimeout int    `toml:"shutdown_timeout"`
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

func (c *Config) validateConfig() error {
	requiredFields := []struct {
		value     string
		fieldName string
	}{
		{c.Server.Host, "server host"},
		{c.Server.JWTSecretKey, "server jwt secret"},
		{c.Database.Host, "database host"},
		{c.Database.Name, "database name"},
		{c.Database.User, "database user"},
		{c.Database.Password, "database password"},
	}

	for _, field := range requiredFields {
		if strings.TrimSpace(field.value) == "" {
			return missingKeyError(field.fieldName)
		}
	}

	if c.Server.Port == 0 {
		return missingKeyError("server port")
	}

	if c.Database.Port == 0 {
		return missingKeyError("database port")
	}

	return nil
}
