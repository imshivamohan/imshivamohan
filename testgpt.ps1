
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
