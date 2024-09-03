# Define paths
$projectPath = "C:\Project" # Adjust this to your project path
$truststorePath = "$projectPath\truststore\truststore.jks"
$dockerComposeFile = "$projectPath\kafka-ui-docker-compose.yml"
$envFilePath = "$projectPath\.env"

# Check if truststore exists
if (-not (Test-Path $truststorePath)) {
    Write-Error "Truststore file does not exist at $truststorePath. Please ensure it is included in the project."
    exit 1
}

# Prompt for keystore path
$keystorePath = Read-Host "Enter the full path to your keystore file"

# Check if the keystore path exists
if (-not (Test-Path $keystorePath)) {
    Write-Error "Keystore path does not exist. Please check the path and try again."
    exit 1
}

# Prompt for other details
$bootstrapServers = Read-Host "Enter the Kafka bootstrap servers (comma-separated)"
$clusterNames = Read-Host "Enter the Kafka cluster names (comma-separated)"

# Set environment variables in .env file
@"
TRUSTSTORE_LOCATION=/etc/kafka/secrets/truststore.jks
TRUSTSTORE_PASSWORD=yourTruststorePassword
KEYSTORE_LOCATION=/etc/kafka/secrets/keystore.jks
KEYSTORE_PASSWORD=yourKeystorePassword
KEY_PASSWORD=yourKeyPassword
BOOTSTRAP_SERVERS=$bootstrapServers
CLUSTER_NAMES=$clusterNames
"@ | Set-Content -Path $envFilePath

# Copy the keystore file to the secrets directory
$secretsPath = "C:\ProgramData\kafka\secrets"
if (-not (Test-Path $secretsPath)) {
    New-Item -Path $secretsPath -ItemType Directory -Force
}

Copy-Item -Path $keystorePath -Destination "$secretsPath\keystore.jks" -Force
Copy-Item -Path $truststorePath -Destination "$secretsPath\truststore.jks" -Force

# Pull and start Docker container
try {
    docker-compose -f $dockerComposeFile pull
    docker-compose -f $dockerComposeFile up -d
} catch {
    Write-Error "Failed to start Docker container. Error: $_"
    exit 1
}

Write-Output "Kafka UI has been installed and is running with the provided keystore and truststore."


#############################


# Define paths and project details
$projectPath = "C:\Project" # Adjust this to your project path
$truststoreSourcePath = "$projectPath\truststore\truststore.jks"
$dockerComposeFile = "$projectPath\kafka-ui-docker-compose.yml"
$envFilePath = "$projectPath\.env"
$secretsPath = "C:\ProgramData\kafka\secrets"

# Ensure the project path and Docker Compose file exist
if (-not (Test-Path $dockerComposeFile)) {
    Write-Error "Docker Compose file does not exist at $dockerComposeFile. Please ensure it is in the correct path."
    exit 1
}

# Prompt for keystore and password details
$keystorePath = Read-Host "Enter the full path to your keystore file"
$keystorePassword = Read-Host "Enter the keystore password" -AsSecureString
$truststorePassword = Read-Host "Enter the truststore password" -AsSecureString
$keystorePasswordPlain = [System.Net.NetworkCredential]::new('', $keystorePassword).Password
$truststorePasswordPlain = [System.Net.NetworkCredential]::new('', $truststorePassword).Password

# Check if the keystore path exists
if (-not (Test-Path $keystorePath)) {
    Write-Error "Keystore path does not exist. Please check the path and try again."
    exit 1
}

# Prompt for Kafka cluster details
$cluster1Name = Read-Host "Enter the name for Cluster 1"
$cluster1BootstrapServers = Read-Host "Enter the bootstrap servers for Cluster 1 (comma-separated)"
$cluster2Name = Read-Host "Enter the name for Cluster 2"
$cluster2BootstrapServers = Read-Host "Enter the bootstrap servers for Cluster 2 (comma-separated)"
$cluster3Name = Read-Host "Enter the name for Cluster 3"
$cluster3BootstrapServers = Read-Host "Enter the bootstrap servers for Cluster 3 (comma-separated)"
$cluster4Name = Read-Host "Enter the name for Cluster 4"
$cluster4BootstrapServers = Read-Host "Enter the bootstrap servers for Cluster 4 (comma-separated)"

# Create the secrets directory if it does not exist
if (-not (Test-Path $secretsPath)) {
    Write-Output "Creating directory $secretsPath"
    New-Item -Path $secretsPath -ItemType Directory -Force
} else {
    Write-Output "Directory $secretsPath already exists"
}

# Copy keystore and truststore files to the secrets directory
Write-Output "Copying keystore file to $secretsPath"
Copy-Item -Path $keystorePath -Destination "$secretsPath\keystore.jks" -Force

Write-Output "Copying truststore file to $secretsPath"
Copy-Item -Path $truststoreSourcePath -Destination "$secretsPath\truststore.jks" -Force

# Create or update the .env file with Kafka cluster configurations
@"
TRUSTSTORE_LOCATION=/etc/kafka/secrets/truststore.jks
TRUSTSTORE_PASSWORD=$truststorePasswordPlain
KEYSTORE_LOCATION=/etc/kafka/secrets/keystore.jks
KEYSTORE_PASSWORD=$keystorePasswordPlain
KEY_PASSWORD=yourKeyPassword

# Kafka Clusters
CLUSTER_1_NAME=$cluster1Name
CLUSTER_1_BOOTSTRAP_SERVERS=$cluster1BootstrapServers
CLUSTER_2_NAME=$cluster2Name
CLUSTER_2_BOOTSTRAP_SERVERS=$cluster2BootstrapServers
CLUSTER_3_NAME=$cluster3Name
CLUSTER_3_BOOTSTRAP_SERVERS=$cluster3BootstrapServers
CLUSTER_4_NAME=$cluster4Name
CLUSTER_4_BOOTSTRAP_SERVERS=$cluster4BootstrapServers
"@ | Set-Content -Path $envFilePath

# Pull and start Docker container with Docker Compose
try {
    Write-Output "Pulling Docker images..."
    docker-compose -f $dockerComposeFile pull
    Write-Output "Starting Docker container..."
    docker-compose -f $dockerComposeFile up -d
} catch {
    Write-Error "Failed to start Docker container. Error: $_"
    exit 1
}

Write-Output "Kafka UI has been installed and is running with the provided keystore and cluster details."
