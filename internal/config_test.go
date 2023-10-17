package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fileContent = `
[database]
host = "localhost"
user = "postgres"
password = "postgres"
port = 5432
name = "postgres"
[server]
port = 8080
host = "localhost"
jwt_secret_key = "xyz"
shutdown_timeout = 10
`
)

// Test read config from string.
func TestReadConfigFromString(t *testing.T) {
	t.Run("read config from string", func(t *testing.T) {
		config, err := ReadConfigFromString(fileContent)
		assert.NoError(t, err)

		expected := Config{
			Database: databaseExpectation(),
			Server:   serverExpectation(),
		}

		assert.Equal(t, expected, config)
	})

	t.Run("Validate missing database host key", func(t *testing.T) {
		// Database host key.
		fileContent := `
[database]
user = "postgres"
password = "postgres"
port = 5432
name = "postgres"
[server]
port = 8080
host = "localhost"
jwt_secret_key = "xyz"
shutdown_timeout = 10
		`
		_, err := ReadConfigFromString(fileContent)
		assert.EqualError(t, err, missingKeyError("database host").Error())
	})

	t.Run("Validate missing database user key", func(t *testing.T) {
		// Database user key.
		fileContent := `
[database]
host = "localhost"
password = "postgres"
port = 5432
name = "postgres"
[server]
port = 8080
host = "localhost"
jwt_secret_key = "xyz"
shutdown_timeout = 10
		`
		_, err := ReadConfigFromString(fileContent)
		assert.EqualError(t, err, missingKeyError("database user").Error())
	})
}

// Test read config from reader.
func TestReadConfigFromReader(t *testing.T) {
	t.Run("read config from reader", func(t *testing.T) {
		reader := strings.NewReader(fileContent)
		config, err := ReadConfigFromReader(reader)
		assert.NoError(t, err)

		expected := Config{
			Database: databaseExpectation(),
			Server:   serverExpectation(),
		}

		assert.Equal(t, expected, config)
	})
}

// The expected database struct, used for testing.
func databaseExpectation() DatabaseConfig {
	return DatabaseConfig{
		Port:     5432,
		Host:     "localhost",
		Name:     "postgres",
		Password: "postgres",
		User:     "postgres",
	}
}

// The expected server struct, used for testing.
func serverExpectation() ServerConfig {
	return ServerConfig{
		Port:            8080,
		Host:            "localhost",
		JWTSecretKey:    "xyz",
		ShutdownTimeout: 10,
	}
}
