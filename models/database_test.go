package models

import (
	"testing"

	internal "github.com/Mahmoud-Emad/envserver/internal"
	"github.com/stretchr/testify/assert"
)

var (
	confContent = `
[database]
host = "localhost"
user = "postgres"
password = "postgres"
port = 5432
name = "postgres"
[server]
port = 8080
host = "localhost"
`
)

// Setup database helper, created to be used inside test case functions.
func setupDB(t *testing.T) (Database, internal.Configuration) {
	db := NewDatabase()
	config, err := internal.ReadConfigFromString(confContent)

	err = db.Connect(config.Database)
	assert.NoError(t, err)

	err = db.Migrate()
	assert.NoError(t, err)

	return db, config
}

// Create a test configuration for database connection.
func createTestConfig() internal.Configuration {
	return internal.Configuration{
		Database: internal.DatabaseConfiguration{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "postgres",
		},
		Server: internal.ServerConfiguration{
			Host: "localhost",
			Port: 8080,
		},
	}
}

// Test connect to database.
// This function expect to initialize a database with a specific configuration.
// The expected behavior is
// 1 - Raising an error in the first block based on the wrong configuration.
// 2 - Passing the second scenario because it's a valid database configuration.
func TestDatabaseConnect(t *testing.T) {
	db, conf := setupDB(t)

	t.Run("invalid database", func(t *testing.T) {
		err := db.Connect(internal.DatabaseConfiguration{})
		assert.Error(t, err)
	})
	t.Run("valid database", func(t *testing.T) {
		err := db.Connect(conf.Database)
		assert.NoError(t, err)
	})
}

func TestUser(t *testing.T) {
	username := "Mahmoud"
	email := "Mahmoud@gmail.com"

	t.Run("create new user object", func(t *testing.T) {
		// Test create new user record into the database.

		db, _ := setupDB(t)
		err := db.CreateUser(&User{
			Name:     username,
			Email:    email,
			Projects: []*Project{},
		})

		assert.NoError(t, err)

		user, err := db.GetUserByEmail(email)

		assert.Equal(t, user.Name, username)
		assert.NoError(t, err)
	})

	t.Run("delete created user", func(t *testing.T) {
		// Test delete user record from the database by it's email.
		db, _ := setupDB(t)
		var user User

		user, err := db.GetUserByEmail(email)

		assert.NoError(t, err)
		assert.Equal(t, user.Name, username)

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
		err := db.CreateProject(&Project{
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

		err := db.CreateProject(&Project{
			Name: projectName,
		})

		assert.NoError(t, err)

		p, err := db.GetProjectByName(projectName)

		// Encrypted value
		encryptedVal, err := internal.EncryptAES([]byte(projectValue), projectKey)
		assert.NoError(t, err)

		err = db.CreateEnvKey(&EnvironmentKey{
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
		encryptedVal, err := internal.EncryptAES([]byte(projectValue), projectKey)
		assert.NoError(t, err)

		decodedVal, err := internal.AESDecryptIt(encryptedVal, projectKey)
		assert.NoError(t, err)

		decodedStoredVal, err := internal.AESDecryptIt(env.Value, projectKey)
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
