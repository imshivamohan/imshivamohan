package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sort"
)

func main() {
	// Define the -f flag for the file path
	filePath := flag.String("f", "", "Path to the security list file")
	flag.Parse()

	// Check if the file path is provided
	if *filePath == "" {
		fmt.Println("Please provide the path to the security list file using -f")
		return
	}

	// Open the specified security.list file
	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create sets to avoid duplicates
	identitiesSet := make(map[string]struct{})
	domainsSet := make(map[string]struct{})
	mapping := make(map[string][]string)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Split the line into domain and identities part
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			fmt.Println("Skipping invalid line:", line)
			continue
		}

		domain := strings.TrimSpace(parts[0])
		identities := strings.Split(strings.TrimSpace(parts[1]), ":")

		// Add the domain to domainsSet
		domainsSet[domain] = struct{}{}

		// Add each identity to identitiesSet and create a mapping
		for _, identity := range identities {
			identity = strings.TrimSpace(identity)
			identitiesSet[identity] = struct{}{}
			mapping[identity] = append(mapping[identity], domain)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Generate SQL statements
	generateSQL(identitiesSet, domainsSet, mapping)
}

func generateSQL(identitiesSet map[string]struct{}, domainsSet map[string]struct{}, mapping map[string][]string) {
	// Generate SQL for technical identities
	fmt.Println("-- List of technical identities")
	var identities []string
	for identity := range identitiesSet {
		identities = append(identities, identity)
	}
	sort.Strings(identities)
	for _, identity := range identities {
		fmt.Printf("INSERT INTO technical_identities (identity) VALUES ('%s');\n", identity)
	}

	// Generate SQL for domains
	fmt.Println("-- List of domains")
	var domains []string
	for domain := range domainsSet {
		domains = append(domains, domain)
	}
	sort.Strings(domains)
	for _, domain := range domains {
		fmt.Printf("INSERT INTO domains (domain_name) VALUES ('%s');\n", domain)
	}

	// Generate SQL for mapping between identities and domains
	fmt.Println("-- Mapping between technical identities and domains")
	for identity, domains := range mapping {
		for _, domain := range domains {
			fmt.Printf("INSERT INTO identity_domain_mapping (identity, domain_name) VALUES ('%s', '%s');\n", identity, domain)
		}
	}
}
