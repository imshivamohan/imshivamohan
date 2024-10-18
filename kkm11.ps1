# Paths to store the encrypted password files
$secretsPath = "C:/ProgramData/Kafka/Secrets/"
$keystorePasswordFile = Join-Path $secretsPath "keystorePassword.xml"
$truststorePasswordFile = Join-Path $secretsPath "truststorePassword.xml"

# Prompt user for passwords
$keystorePassword = Read-Host "Enter keystore password" -AsSecureString
$truststorePassword = Read-Host "Enter truststore password" -AsSecureString

# Store passwords in an encrypted format
$keystorePassword | Export-CliXml -Path $keystorePasswordFile
$truststorePassword | Export-CliXml -Path $truststorePasswordFile

Write-Host "Passwords have been securely stored."



###################################################################################

# Define paths
$secretsPath = "C:/ProgramData/Kafka/Secrets/"
$keystoreFile = Join-Path $secretsPath "keystore.jks"
$truststoreFile = Join-Path $secretsPath "truststore.jks"
$keystorePasswordFile = Join-Path $secretsPath "keystorePassword.xml"
$truststorePasswordFile = Join-Path $secretsPath "truststorePassword.xml"
$envFile = ".env"
$dockerComposeFile = "docker-compose.yml"  # Replace with actual path to your docker-compose.yml

# Function to check if a process is running
function Check-ServiceStatus {
    param ($serviceName)
    $process = Get-Process -Name $serviceName -ErrorAction SilentlyContinue
    return $process -ne $null
}

# Check if Rancher Desktop and Docker are running
$dockerRunning = Check-ServiceStatus -serviceName "com.docker.service"
$rancherDesktopRunning = Check-ServiceStatus -serviceName "Rancher Desktop"

if ($dockerRunning -and $rancherDesktopRunning) {
    Write-Host "Rancher Desktop and Docker are both running."

    # Check for keystore and truststore files
    if (-Not (Test-Path $keystoreFile) -or -Not (Test-Path $truststoreFile)) {
        Write-Host "Keystore or Truststore file missing."

        # Prompt user for the keystore and truststore files
        $keystorePath = Read-Host "Please enter the path to the keystore.jks file"
        $truststorePath = Read-Host "Please enter the path to the truststore.jks file"

        # Copy the files to the secrets directory
        Copy-Item -Path $keystorePath -Destination $keystoreFile
        Copy-Item -Path $truststorePath -Destination $truststoreFile

        Write-Host "Keystore and Truststore files copied to $secretsPath."
    } else {
        Write-Host "Keystore and Truststore files are already available."
    }

    # Check if password files exist and read them
    if (-Not (Test-Path $keystorePasswordFile) -or -Not (Test-Path $truststorePasswordFile)) {
        Write-Host "Password files are missing. Please store them first."
        exit
    }

    # Read the passwords from the encrypted files
    $keystorePassword = Import-CliXml -Path $keystorePasswordFile
    $truststorePassword = Import-CliXml -Path $truststorePasswordFile

    # Create .env file for Kafka cluster configuration
    $kafkaBootstrapServers = Read-Host "Enter Kafka Bootstrap Servers"
    $kafkaClusterName = Read-Host "Enter Kafka Cluster Name"

    $envContent = @"
KAFKA_BOOTSTRAP_SERVERS=$kafkaBootstrapServers
KAFKA_CLUSTER_NAME=$kafkaClusterName
KEYSTORE_PASSWORD=$(ConvertFrom-SecureString $keystorePassword -AsPlainText)
TRUSTSTORE_PASSWORD=$(ConvertFrom-SecureString $truststorePassword -AsPlainText)
"@
    Set-Content -Path $envFile -Value $envContent

    Write-Host ".env file created with Kafka configuration."

    # Pull and start Docker container using Docker Compose
    Write-Host "Pulling and starting Docker container..."
    docker-compose -f $dockerComposeFile up -d --build
    Write-Host "Docker container started successfully."
} else {
    Write-Host "Rancher Desktop or Docker is not running. Please start both services and try again."
}
