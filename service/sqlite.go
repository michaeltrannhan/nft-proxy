package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	nft_proxy "github.com/alphabatem/nft-proxy"
	"github.com/babilu-online/common/context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DBError represents a database error with HTTP status code
type DBError struct {
	StatusCode int
	Err        error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("database error (status %d): %v", e.StatusCode, e.Err)
}

type SqliteService struct {
	context.DefaultService
	db *gorm.DB

	username string
	password string
	database string
	host     string
}

const SQLITE_SVC = "sqlite_svc"

// Id returns Service ID
func (ds SqliteService) Id() string {
	return SQLITE_SVC
}

// Db Access to raw SqliteService db
func (ds SqliteService) Db() *gorm.DB {
	return ds.db
}

// Configure the service
func (ds *SqliteService) Configure(ctx *context.Context) error {
	ds.database = os.Getenv("DB_DATABASE")

	return ds.DefaultService.Configure(ctx)
}

// Start the service and open connection to the database
// Migrate any tables that have changed since last runtime
func (ds *SqliteService) Start() (err error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	}

	ds.db, err = gorm.Open(sqlite.Open(ds.database), config)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := ds.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Add indexes and migrate schema
	err = ds.db.AutoMigrate(&nft_proxy.SolanaMedia{})
	if err != nil {
		return fmt.Errorf("failed to migrate schema: %w", err)
	}

	return nil
}

// Find returns the db query for a statement
func (ds *SqliteService) Find(out interface{}, where string, args ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ds.error(ds.db.WithContext(ctx).Find(out, where, args).Error)
}

// Create a new item in the SqliteService
func (ds *SqliteService) Create(val interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ds.db.WithContext(ctx).Create(val).Error
	return val, ds.error(err)
}

// Update an existing item
func (ds *SqliteService) Update(old interface{}, new interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ds.db.WithContext(ctx).Model(old).Updates(new).Error
	return new, ds.error(err)
}

// Delete an existing item
func (ds *SqliteService) Delete(val interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return ds.error(ds.db.WithContext(ctx).Delete(val).Error)
}

// Migrate creates any new tables needed
func (ds *SqliteService) Migrate(values ...interface{}) error {
	err := ds.db.AutoMigrate(values...)
	if err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

// Shutdown Gracefully close the database connection
func (ds *SqliteService) Shutdown() {
	// Close the database connection
	sqlDB, err := ds.db.DB()
	if err != nil {
		log.Println(err)
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Println(err)
	}
}

// Parse an error returned from the database into a more contextual error that can be used with http response codes
func (ds *SqliteService) error(err error) error {
	if err == nil {
		return nil
	}

	var statusCode int

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		statusCode = 404
	case errors.Is(err, context.DeadlineExceeded):
		statusCode = 504
	default:
		statusCode = 500
	}

	log.Printf("Database error: %v (status code: %d)", err, statusCode)
	return &DBError{
		StatusCode: statusCode,
		Err:        err,
	}
}
