package dbcon

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus" // Logging library
	_ "github.com/lib/pq"        // PostgreSQL driver
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	_ "modernc.org/sqlite"       // SQLite driver
)

var logger = logrus.New()

// DBWrapper wraps the database connection and provides additional functionality.
type DBWrapper struct {
	DB *sql.DB
}

// Config represents database configuration.
type Config struct {
	Driver   string `yaml:"driver"`   // Database driver (e.g., "postgres", "mysql", "sqlite3")
	Host     string `yaml:"host"`     // Hostname or IP
	Port     int    `yaml:"port"`     // Port number
	User     string `yaml:"user"`     // Username
	Password string `yaml:"password"` // Password
	DBName   string `yaml:"dbname"`   // Database name
	SSLMode  string `yaml:"sslmode"`  // SSL mode (for PostgreSQL)
}

// NewDB initializes a database connection using the provided config.
func NewDB(cfg Config) (*DBWrapper, error) {
	var dsn string

	switch cfg.Driver {
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	case "sqlite3":
		dsn = cfg.DBName // SQLite uses only the database file path as DSN
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		logger.Error("Failed to connect to database: ", err)
		return nil, err
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		logger.Error("Failed to ping database: ", err)
		return nil, err
	}

	logger.Info("Database connection established")
	return &DBWrapper{DB: db}, nil
}

// Query executes a SELECT query.
func (dbw *DBWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := dbw.DB.Query(query, args...)
	elapsed := time.Since(start)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
			"error": err,
		}).Error("Query failed")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"query":   query,
		"args":    args,
		"elapsed": elapsed,
	}).Info("Query executed")
	return rows, nil
}

// Exec executes an INSERT, UPDATE, or DELETE query.
func (dbw *DBWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := dbw.DB.Exec(query, args...)
	elapsed := time.Since(start)

	if err != nil {
		logger.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
			"error": err,
		}).Error("Exec failed")
		return nil, err
	}

	logger.WithFields(logrus.Fields{
		"query":   query,
		"args":    args,
		"elapsed": elapsed,
	}).Info("Exec executed")
	return result, nil
}

// Close closes the database connection.
func (dbw *DBWrapper) Close() error {
	logger.Info("Closing database connection")
	return dbw.DB.Close()
}
######################################
Step 4: YAML Configuration File Example
Create a config.yaml for database connection settings.

yaml
Copy code
driver: "postgres"
host: "localhost"
port: 5432
user: "username"
password: "password"
dbname: "example_db"
sslmode: "disable"


#####################

Step 5: Integrate the Library
Main Application Example (main.go)
go
Copy code
package main

import (
	"fmt"
	"log"

	"your_project/dbcon" // Replace with your library's import path

	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	// Read database config from YAML
	configFile := "config.yaml"
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg dbcon.Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Initialize the database
	db, err := dbcon.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Example query
	query := "SELECT id, name FROM users WHERE age > $1"
	rows, err := db.Query(query, 18)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("User: %d, %s\n", id, name)
	}
}

$##################################
Features:
Support for Multiple Databases: Easily switch between PostgreSQL, MySQL, and SQLite by modifying the YAML configuration.
Logging: Provides structured logging for queries and execution times.
Reusability: Encapsulates common database operations in a reusable library.
Let me know if you need further guidance


