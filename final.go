Go's database/sql package has a built-in connection pooling mechanism that manages the connection pool automatically. However, you can configure its behavior to optimize performance based on the application's needs.

Here's how you can configure the connection pool in dbcon/db.go by setting properties like maximum open connections, maximum idle connections, and connection lifetime.

Adding Connection Pool Settings in dbcon/db.go
Let's update the Connect function to include connection pool configuration.

go
Copy code
package dbcon

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"            // PostgreSQL driver
    _ "github.com/go-sql-driver/mysql" // MySQL driver
    _ "github.com/godror/godror"      // Oracle driver
    "github.com/spf13/viper"
    "log"
    "time"
)

// DB struct to hold the database connection
type DB struct {
    Conn *sql.DB
}

// InitConfig initializes the configuration using viper
func InitConfig(configFilePath string) error {
    viper.SetConfigFile(configFilePath)
    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("error reading config file: %w", err)
    }
    return nil
}

// Connect initializes a database connection with connection pool settings
func Connect() (*DB, error) {
    driver := viper.GetString("database.driver")
    host := viper.GetString("database.host")
    port := viper.GetInt("database.port")
    user := viper.GetString("database.user")
    password := viper.GetString("database.password")
    dbname := viper.GetString("database.dbname")
    service := viper.GetString("database.service")

    var dsn string
    switch driver {
    case "postgres":
        dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
            host, port, user, password, dbname)
    case "mysql":
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
            user, password, host, port, dbname)
    case "oracle":
        dsn = fmt.Sprintf("%s/%s@%s:%d/%s", user, password, host, port, service)
    default:
        return nil, fmt.Errorf("unsupported database driver: %s", driver)
    }

    db, err := sql.Open(driver, dsn)
    if err != nil {
        return nil, err
    }

    // Set up connection pool configuration
    maxOpenConns := viper.GetInt("database.max_open_conns")
    maxIdleConns := viper.GetInt("database.max_idle_conns")
    connMaxLifetime := viper.GetDuration("database.conn_max_lifetime")

    db.SetMaxOpenConns(maxOpenConns)
    db.SetMaxIdleConns(maxIdleConns)
    db.SetConnMaxLifetime(connMaxLifetime)

    // Verify connection
    if err := db.Ping(); err != nil {
        db.Close()
        return nil, err
    }

    log.Println("Database connected successfully with connection pooling")
    return &DB{Conn: db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
    return db.Conn.Close()
}
Updating Configuration File (config.yaml)
You can now add connection pooling settings in config.yaml. Here’s an example:

yaml
Copy code
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  user: "your_user"
  password: "your_password"
  dbname: "your_database"
  service: "your_service" # Oracle-specific field

  # Connection pooling settings
  max_open_conns: 25             # Max open connections to the database
  max_idle_conns: 10             # Max idle connections to maintain in the pool
  conn_max_lifetime: 30m         # Maximum amount of time a connection can live (30 minutes here)
Explanation of Connection Pool Parameters
SetMaxOpenConns: This sets the maximum number of open connections to the database. Setting this to a value greater than 0 limits the total number of open connections, which can improve performance under heavy load by reducing contention.

SetMaxIdleConns: This controls the maximum number of connections in the idle connection pool. If there are more idle connections than specified by this value, excess connections are closed. Setting this appropriately can reduce the overhead of opening new connections.

SetConnMaxLifetime: This defines the maximum duration a connection can be reused before it is closed. This is helpful for managing resources and ensuring that stale connections (or connections held open for too long) are closed and recycled.

Usage in main.go
This usage example remains the same:

go
Copy code
package main

import (
    "fmt"
    "log"
    "your_project_path/dbcon"
)

func main() {
    configPath := "path/to/config.yaml"

    // Initialize configuration
    if err := dbcon.InitConfig(configPath); err != nil {
        log.Fatalf("Failed to initialize configuration: %v", err)
    }

    // Connect to the database
    db, err := dbcon.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    fmt.Println("Database connection established!")
}
Summary
This setup enables you to configure connection pooling through the YAML configuration file, making it easy to adjust pool settings without modifying the code.






You said:
above code is looks good how to do package this module in gitlab.com my gitlab url will be "https://devcloude.ubs.net/ubs/eis/golang-common/" i have toke name as "Eis-token" and toke is "adfdsfasd3q5345"
ChatGPT said:
ChatGPT
To package and host your Golang database connection library on GitLab, follow these steps:

Create a New GitLab Repository

Visit your GitLab instance at https://devcloude.ubs.net/ubs/eis/golang-common/.
Create a new repository with the name golang-common.
Add a .gitignore file for Go if available.
Structure the Go Module

Organize your code under a project folder with a go.mod file, which is necessary for Go modules.
Initialize Your Go Module

Open a terminal in your project directory and initialize the module:

sh
Copy code
go mod init devcloude.ubs.net/ubs/eis/golang-common
This command sets up go.mod, where devcloude.ubs.net/ubs/eis/golang-common becomes the module path.

Add the Code and Commit Changes

Ensure your folder structure is set up like this:

go
Copy code
golang-common/
├── dbcon/
│   ├── db.go
├── config.yaml
├── go.mod
├── go.sum
└── README.md
Then, add your files to Git:

sh
Copy code
git init
git remote add origin https://devcloude.ubs.net/ubs/eis/golang-common.git
git add .
git commit -m "Initial commit for common Go database connection library"
git push -u origin main
Add GitLab Authentication for Private Repositories

To use this repository as a dependency in other projects, set up GitLab authentication. Since the repository is private, use a personal access token (PAT) for authentication.

Token Setup: Ensure Eis-token (PAT) has permissions for read access to this repository.
Environment Setup: Add the PAT to your Git configuration to authenticate GitLab pulls.
Configure Git to cache credentials:

sh
Copy code
git config --global credential.helper cache
Add your GitLab token using the following command:

sh
Copy code
git config --global url."https://Eis-token:adfdsfasd3q5345@devcloude.ubs.net/ubs/eis/".insteadOf "https://devcloude.ubs.net/ubs/eis/"
Using the Module in Other Projects

To use this module in another Go project, add it as a dependency by running:

sh
Copy code
go get devcloude.ubs.net/ubs/eis/golang-common@latest
Go will now use the credentials set in step 5 to access the module from GitLab.

Example Usage in a Project
In another Go project, import the dbcon package:

go
Copy code
import "devcloude.ubs.net/ubs/eis/golang-common/dbcon"
Then, follow the normal steps to initialize and use the connection as shown in previous examples.

This setup allows you to keep the code private while still enabling access via Go’s module system.
