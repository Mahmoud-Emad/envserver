package internal

import (
	"fmt"
	"reflect"

	models "github.com/Mahmoud-Emad/envserver/models"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB struct hold db instance
type Database struct {
	db *gorm.DB
}

// NewDatabase create and return new Database struct.
func NewDatabase() Database {
	return Database{}
}

// Connect connects to database server.
func (d *Database) Connect(dbConfig DatabaseConfiguration) error {
	log.Info().Msg("Connecting to the database.")
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		return err
	}

	d.db = gormDB
	log.Info().Msg("Database Connected.")
	return nil
}

// Migrate migrates the database schema.
func (d *Database) Migrate() error {
	tables := []interface{}{&models.User{}, &models.Project{}, &models.EnvironmentKey{}}

	log.Info().Msg("Database migration started")
	for _, table := range tables {
		tableName := reflect.TypeOf(table).Elem().Name()
		log.Info().Msgf("Migrating table: %s", tableName)
		if err := d.db.AutoMigrate(table); err != nil {
			log.Error().Msgf("failed to migrate table %s: %q", tableName, err)
			return err
		}
	}

	log.Info().Msg("Database migration completed")
	return nil
}

// Create new user object inside the daabase.
func (d *Database) CreateUser(u *models.User) error {
	result := d.db.Create(&u)
	return result.Error
}

// GetUserByEmail returns user by its email
func (d *Database) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	query := d.db.First(&u, "email = ?", email)
	return u, query.Error
}

// GetUserByID returns user by its email
func (d *Database) GetUserByID(id uint) (models.User, error) {
	var u models.User
	query := d.db.First(&u, "id = ?", id)
	return u, query.Error
}

// GetProjectByID returns user by its email
func (d *Database) GetProjectByID(id uint) (models.Project, error) {
	var p models.Project
	query := d.db.First(&p, "id = ?", id)
	return p, query.Error
}

// GetUsers returns a list of all user records
func (d *Database) GetUsers() ([]models.User, error) {
	// Retrieve all users
	var users []models.User
	result := d.db.Find(&users)
	return users, result.Error
}

// GetUsers returns a list of all project records
func (d *Database) GetProjects() ([]models.Project, error) {
	// Retrieve all users
	var projects []models.Project
	result := d.db.Find(&projects)
	return projects, result.Error
}

// DeleteUserByEmail deletes a user by their email
func (d *Database) DeleteUserByEmail(email string) error {
	result := d.db.Unscoped().Where("email = ?", email).Delete(&models.User{})
	return result.Error
}

// DeleteUserByID deletes a user by their id
func (d *Database) DeleteUserByID(id uint) error {
	result := d.db.Unscoped().Where("id = ?", id).Delete(&models.User{})
	return result.Error
}

// Create new project object inside the daabase.
func (d *Database) CreateProject(p *models.Project) error {
	result := d.db.Create(&p)
	return result.Error
}

// DeleteProjectByName deletes a project by it's name
func (d *Database) DeleteProjectByName(name string) error {
	result := d.db.Unscoped().Where("name = ?", name).Delete(&models.Project{})
	return result.Error
}

// DeleteProjectByName deletes a project by it's name
func (d *Database) DeleteProjectByID(id uint) error {
	result := d.db.Unscoped().Where("id = ?", id).Delete(&models.Project{})
	return result.Error
}

// GetProjectByName returns user by its name
func (d *Database) GetProjectByName(name string) (models.Project, error) {
	var p models.Project
	query := d.db.First(&p, "name = ?", name)
	return p, query.Error
}

// Create new EnvironmentKey object inside the daabase.
func (d *Database) CreateEnvKey(env *models.EnvironmentKey) error {
	result := d.db.Create(&env)
	return result.Error
}

// DeleteEnvironmentKeyByKeyName Delete an EnvironmentKey by it's key name.
func (d *Database) DeleteEnvKeyByKeyName(keyName string) error {
	result := d.db.Unscoped().Where("key = ?", keyName).Delete(&models.EnvironmentKey{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetEnvKeyByKeyName returns user by its key name
func (d *Database) GetEnvKeyByKeyName(keyName string) (models.EnvironmentKey, error) {
	var env models.EnvironmentKey
	query := d.db.First(&env, "key = ?", keyName)
	return env, query.Error
}