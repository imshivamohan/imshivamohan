1. Set Up the GORM Database Library
We will create a db package with the following structure:

db/connection.go: For database connection setup and configuration.
db/logger.go: For logging support.
db/models.go: Optional file for defining models for your database schema.
1.1 db/connection.go
This file will handle the database connection setup, including pooling and configuration.

go
Copy code
package db

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Config holds database configuration details
type Config struct {
	Driver          string
	User            string
	Password        string
	Host            string
	Port            string
	Database        string
	SSLMode         string
	PoolMaxOpenConns int
	PoolMaxIdleConns int
	PoolMaxLifetime  time.Duration
	LogLevel         string
}

// Initialize connects to the database with the provided configuration
func Initialize(cfg Config) error {
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
		)
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
		)
		dialector = mysql.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(cfg.Database) // For SQLite, the "Database" field is the file path
	default:
		return fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	// Configure logging
	var gormLogger logger.Interface
	switch cfg.LogLevel {
	case "silent":
		gormLogger = logger.Default.LogMode(logger.Silent)
	case "error":
		gormLogger = logger.Default.LogMode(logger.Error)
	case "warn":
		gormLogger = logger.Default.LogMode(logger.Warn)
	case "info":
		gormLogger = logger.Default.LogMode(logger.Info)
	default:
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Initialize GORM
	var err error
	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	// Configure connection pooling
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to configure connection pool: %v", err)
	}
	sqlDB.SetMaxOpenConns(cfg.PoolMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.PoolMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.PoolMaxLifetime)

	log.Println("Database connection established with pooling.")
	return nil
}

// Close closes the database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to retrieve SQL DB instance: %v", err)
	}
	return sqlDB.Close()
}
1.2 db/logger.go
This file provides utility functions for custom logging.

go
Copy code
package db

import (
	"log"
)

var logLevel string

// SetLogLevel sets the logging level
func SetLogLevel(level string) {
	logLevel = level
}

// Debug logs debug messages
func Debug(msg string) {
	if logLevel == "debug" {
		log.Println("[DEBUG]", msg)
	}
}

// Info logs informational messages
func Info(msg string) {
	if logLevel == "info" || logLevel == "debug" {
		log.Println("[INFO]", msg)
	}
}

// Error logs error messages
func Error(msg string) {
	log.Println("[ERROR]", msg)
}
2. Usage Examples
Here are examples showing how to pass database configurations and use the library.

2.1 Passing Configuration with Code
go
Copy code
package main

import (
	"log"
	"time"
	"your_project_path/db"
)

func main() {
	// Define database configuration
	cfg := db.Config{
		Driver:          "postgres",
		User:            "your_user",
		Password:        "your_password",
		Host:            "localhost",
		Port:            "5432",
		Database:        "your_database",
		SSLMode:         "disable",
		PoolMaxOpenConns: 25,
		PoolMaxIdleConns: 10,
		PoolMaxLifetime:  5 * time.Minute,
		LogLevel:         "info",
	}

	// Initialize the database
	err := db.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Example query: Create a record
	type User struct {
		ID   uint
		Name string
	}
	err = db.DB.Create(&User{Name: "John Doe"}).Error
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	log.Println("User created successfully.")
}
2.2 Passing Configuration from a File
Using a YAML file for configuration:

config.yaml
yaml
Copy code
driver: postgres
user: your_user
password: your_password
host: localhost
port: 5432
database: your_database
sslmode: disable
pool_max_open_conns: 25
pool_max_idle_conns: 10
pool_max_lifetime: 300
log_level: info
Code to Read YAML and Initialize:
go
Copy code
package main

import (
	"log"
	"time"
	"your_project_path/db"

	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	// Read configuration from YAML file
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	var cfg db.Config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	// Convert PoolMaxLifetime to duration
	cfg.PoolMaxLifetime = time.Duration(cfg.PoolMaxLifetime) * time.Second

	// Initialize the database
	err = db.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Example query
	type User struct {
		ID   uint
		Name string
	}
	err = db.DB.Create(&User{Name: "Jane Doe"}).Error
	if err != nil {
		log.Fatalf("Error creating user: %v", err)
	}
	log.Println("User created successfully.")
}
2.3 Passing Environment Variables
Set environment variables:

bash
Copy code
export DB_DRIVER=postgres
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=your_database
export DB_SSLMODE=disable
export DB_POOL_MAX_OPEN_CONNS=25
export DB_POOL_MAX_IDLE_CONNS=10
export DB_POOL_MAX_LIFETIME=300
export DB_LOG_LEVEL=info
Load them dynamically in your program using os.Getenv.

3. Supported Drivers
PostgreSQL: "postgres"
MySQL: "mysql"
SQLite: "sqlite"
You can extend the library to support other drivers by adding their configurations in the Initialize function.

This setup provides a modular, reusable, and flexible database connection library with GORM.
