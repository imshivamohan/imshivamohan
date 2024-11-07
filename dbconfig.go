
database:
  driver: "postgres"  # or "mysql", "sqlite", etc.
  host: "localhost"
  port: 5432
  user: "username"
  password: "password"
  dbname: "your_database"

####################################


package dbutils

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}

// LoadConfig loads the database configuration from a YAML file.
func LoadConfig(file string) (*Config, error) {
	config := &Config{}
	yamlFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

###########################################

package dbutils

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Import the required database driver
)

// DB is a wrapper struct that holds the SQL database connection.
type DB struct {
	*sql.DB
}

// ConnectDB establishes a connection to the database using the provided configuration.
func ConnectDB(cfg *Config) (*DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.DBName)
	db, err := sql.Open(cfg.Database.Driver, connStr)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

// Insert executes an INSERT query with given parameters.
func (db *DB) Insert(query string, args ...interface{}) error {
	_, err := db.Exec(query, args...)
	return err
}


############################################


package main

import (
	"flag"
	"log"

	"path/to/your/dbutils" // Update this with your actual package path
)

func main() {
	// Get config file path from command-line argument
	configFile := flag.String("f", "config.yaml", "Path to the config file")
	flag.Parse()

	// Load database configuration
	cfg, err := dbutils.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the database
	db, err := dbutils.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Insert data example
	query := "INSERT INTO your_table (column1, column2) VALUES ($1, $2)"
	err = db.Insert(query, "value1", "value2")
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}

	log.Println("Data inserted successfully")
}


############################################


go run main.go -f config.yaml
