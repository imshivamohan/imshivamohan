-- Create table for technical identities
CREATE TABLE technical_identities (
    identity VARCHAR(255) UNIQUE NOT NULL, -- Enforces unique identities and prevents nulls
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically stores the insert timestamp
    PRIMARY KEY (identity)
);

-- Create table for data domains
CREATE TABLE data_domains (
    domain_name VARCHAR(255) UNIQUE NOT NULL, -- Enforces unique domains and prevents nulls
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically stores the insert timestamp
    PRIMARY KEY (domain_name)
);

-- Create table for mapping between technical identities and data domains
CREATE TABLE data_domain_identities (
    identity VARCHAR(255) NOT NULL,
    domain_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Automatically stores the insert timestamp
    PRIMARY KEY (identity, domain_name), -- Prevents duplicate identity-domain mappings
    FOREIGN KEY (identity) REFERENCES technical_identities(identity),
    FOREIGN KEY (domain_name) REFERENCES data_domains(domain_name)
);

package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	filePath := flag.String("f", "", "Path to the security.list file")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please provide a file path with the -f flag.")
		return
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	identities := make(map[string]struct{})
	domains := make(map[string]struct{})
	mappings := make(map[string]map[string]struct{})

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	for _, record := range records {
		if len(record) != 2 {
			continue
		}
		domain := strings.TrimSpace(record[0])
		identityList := strings.Split(strings.TrimSpace(record[1]), ":")

		domains[domain] = struct{}{}

		if _, exists := mappings[domain]; !exists {
			mappings[domain] = make(map[string]struct{})
		}

		for _, identity := range identityList {
			identity = strings.TrimSpace(identity)
			identities[identity] = struct{}{}
			mappings[domain][identity] = struct{}{}
		}
	}

	connStr := "user=yourusername dbname=yourdbname sslmode=disable" 
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	for identity := range identities {
		_, err := db.Exec("INSERT INTO technical_identities (identity) VALUES ($1) ON CONFLICT (identity) DO NOTHING", identity)
		if err != nil {
			fmt.Println("Error inserting identity:", err)
		}
	}

	for domain := range domains {
		_, err := db.Exec("INSERT INTO data_domains (domain_name) VALUES ($1) ON CONFLICT (domain_name) DO NOTHING", domain)
		if err != nil {
			fmt.Println("Error inserting domain:", err)
		}
	}

	for domain, identitySet := range mappings {
		for identity := range identitySet {
			_, err := db.Exec("INSERT INTO data_domain_identities (identity, domain_name) VALUES ($1, $2) ON CONFLICT (identity, domain_name) DO NOTHING", identity, domain)
			if err != nil {
				fmt.Println("Error inserting mapping:", err)
			}
		}
	}

	fmt.Println("Data inserted successfully!")
}
