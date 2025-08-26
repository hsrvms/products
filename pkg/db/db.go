package db

import (
	"fmt"
	"go-starter/pkg/config"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB  *gorm.DB
	cfg *config.Config
}

func New(cfg *config.Config) (*Database, error) {
	dsn := cfg.GetDBConnString()

	// Configure GORM logger based on environment
	logLevel := logger.Info
	if os.Getenv("APP_ENV") == "production" {
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Enhanced connection pool settings
	maxOpenConns := 25
	if os.Getenv("APP_ENV") == "production" {
		maxOpenConns = 50
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	// Test the connection with timeout
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully with %d max connections", maxOpenConns)

	log.Println("Successfully connected to the database with GORM")
	return &Database{DB: db, cfg: cfg}, nil
}

func (d *Database) Close() error {
	if d.DB != nil {
		sqlDB, err := d.DB.DB()
		if err != nil {
			return err
		}
		if err := sqlDB.Close(); err != nil {
			return err
		}
		log.Println("Database connection closed")
	}
	return nil
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}

// MigrateWithCleanup handles migration from existing schema to GORM schema
func (d *Database) MigrateWithCleanup(models ...interface{}) error {
	// First, check if the table exists and clean up old constraints
	var tableExists bool
	err := d.DB.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists).Error
	if err != nil {
		return fmt.Errorf("failed to check table existence: %w", err)
	}

	if tableExists {
		log.Println("Existing users table found, cleaning up constraints...")

		// Drop old constraints that might conflict with GORM
		cleanupSQL := []string{
			"ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key",
			"DROP INDEX IF EXISTS idx_users_email",
		}

		for _, sql := range cleanupSQL {
			if err := d.DB.Exec(sql).Error; err != nil {
				log.Printf("Warning: Could not execute cleanup SQL '%s': %v", sql, err)
			}
		}
	}

	// Now run GORM auto-migration
	return d.DB.AutoMigrate(models...)
}
