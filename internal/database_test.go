package internal

import (
	"testing"

	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/stretchr/testify/assert"
)

var configContent = `
[server]
host = "localhost"
port = 8080
jwt_secret_key = "xyz"
shutdown_timeout = 10

[database]
host = "localhost"
port = 5432
name = "postgres"
user = "postgres"
password = "postgres"
`

// Setup database helper, created to be used inside test case functions.
func setupDB(t *testing.T) (Database, Config) {
	db := NewDatabase()
	config, err := ReadConfigFromString(configContent)

	err = db.Connect(config.Database)
	assert.NoError(t, err)

	err = db.Migrate()
	assert.NoError(t, err)

	return db, config
}

// Create a test Config for database connection.
func createTestConfig() Config {
	return Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "postgres",
		},

		Server: ServerConfig{
			Host: "localhost",
			Port: 8080,
		},
	}
}

// Test connect to database.
// This function expect to initialize a database with a specific Config.
// The expected behavior is
// 1 - Raising an error in the first block based on the wrong Config.
// 2 - Passing the second scenario because it's a valid database Config.
func TestDatabaseConnect(t *testing.T) {
	db, conf := setupDB(t)

	t.Run("invalid database", func(t *testing.T) {
		err := db.Connect(DatabaseConfig{})
		assert.Error(t, err)
	})

	t.Run("valid database", func(t *testing.T) {
		err := db.Connect(conf.Database)
		assert.NoError(t, err)
	})
}

func TestUser(t *testing.T) {
	FName := "Mahmoud"
	LName := "Emad"
	email := "Mahmoud@gmail.com"

	t.Run("create new user object", func(t *testing.T) {
		// Test create new user record into the database.

		db, _ := setupDB(t)
		err := db.CreateUser(&models.User{
			FirstName: FName,
			LastName:  LName,
			Email:     email,
			Projects:  []*models.Project{},
		})

		assert.NoError(t, err)
		user, err := db.GetUserByEmail(email)

		assert.Equal(t, user.FirstName, FName)
		assert.Equal(t, user.LastName, LName)
		assert.NoError(t, err)
	})

	t.Run("delete created user", func(t *testing.T) {
		// Test delete user record from the database by it's email.
		db, _ := setupDB(t)
		var user models.User
		user, err := db.GetUserByEmail(email)

		assert.NoError(t, err)
		assert.Equal(t, user.FirstName, FName)
		assert.Equal(t, user.LastName, LName)

		err = db.DeleteUserByEmail(email)
		assert.NoError(t, err)

		user, err = db.GetUserByEmail(email)
		assert.Error(t, err)
	})
}

func TestProject(t *testing.T) {
	projectName := "ligdude"

	t.Run("create new project object", func(t *testing.T) {
		// Test create new project record into the database.
		db, _ := setupDB(t)
		err := db.CreateProject(&models.Project{
			Name: projectName,
		})

		assert.NoError(t, err)

		p, err := db.GetProjectByName(projectName)

		assert.Equal(t, p.Name, projectName)
		assert.NoError(t, err)
	})

	t.Run("delete created project", func(t *testing.T) {
		// Test delete project record from the database by its name.
		db, _ := setupDB(t)

		p, err := db.GetProjectByName(projectName)

		assert.NoError(t, err)
		assert.Equal(t, p.Name, projectName)

		err = db.DeleteProjectByName(projectName)
		assert.NoError(t, err)

		_, err = db.GetProjectByName(projectName)
		assert.Error(t, err)
	})
}

func TestEnvironmentKey(t *testing.T) {
	projectName := "ligdude"

	// Config key/value
	projectKey := "password"
	projectValue := "xyz@M@#Jois2$#!"

	t.Run("create new environment Key object", func(t *testing.T) {
		// Test create new env key|value record into the database.
		db, _ := setupDB(t)

		err := db.CreateProject(&models.Project{
			Name: projectName,
		})

		assert.NoError(t, err)
		p, err := db.GetProjectByName(projectName)

		// Encrypted value
		encryptedVal, err := EncryptAES([]byte(projectValue), projectKey)
		assert.NoError(t, err)

		err = db.CreateEnvKey(&models.EnvironmentKey{
			Key:       projectKey,
			Value:     encryptedVal,
			ProjectID: p.ID,
		})
		assert.NoError(t, err)

		env, err := db.GetEnvKeyByKeyName(projectKey)
		assert.Equal(t, env.Key, projectKey)
		assert.NoError(t, err)
	})

	t.Run("delete created environment key", func(t *testing.T) {
		// Test delete key|value record from the database by its key name.
		db, _ := setupDB(t)

		env, err := db.GetEnvKeyByKeyName(projectKey)
		assert.NoError(t, err)

		// Encrypted value
		encryptedVal, err := EncryptAES([]byte(projectValue), projectKey)
		assert.NoError(t, err)

		decodedVal, err := DecryptAES(encryptedVal, projectKey)
		assert.NoError(t, err)

		decodedStoredVal, err := DecryptAES(env.Value, projectKey)
		assert.NoError(t, err)

		assert.Equal(t, decodedVal, decodedStoredVal)
		assert.Equal(t, string(decodedVal), string(decodedStoredVal))

		err = db.DeleteEnvKeyByKeyName(projectKey)
		assert.NoError(t, err)

		_, err = db.GetEnvKeyByKeyName(projectKey)
		assert.Error(t, err)

		err = db.DeleteProjectByName(projectName)
		assert.NoError(t, err)

		_, err = db.GetProjectByName(projectName)
		assert.Error(t, err)
	})
}
