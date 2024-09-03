# Prompt for keystore and truststore locations
$keystorePath = Read-Host "Enter the path to the keystore file (e.g., C:\path\to\keystore.jks)"
$truststorePath = Read-Host "Enter the path to the truststore file (e.g., C:\path\to\truststore.jks)"
$keystorePassword = Read-Host "Enter the keystore password" -AsSecureString
$truststorePassword = Read-Host "Enter the truststore password" -AsSecureString

# Convert SecureString to plain text for further use
$keystorePasswordPlain = [System.Net.NetworkCredential]::new($keystorePassword).Password
$truststorePasswordPlain = [System.Net.NetworkCredential]::new($truststorePassword).Password

# Install Docker
choco install docker-desktop

# Create directory for SSL files if it does not exist
$secretsPath = "C:\ProgramData\kafka\secrets"
if (-not (Test-Path $secretsPath)) {
    New-Item -Path $secretsPath -ItemType Directory
}

# Copy the keystore and truststore files to the secrets directory
Copy-Item -Path $keystorePath -Destination "$secretsPath\kafka.keystore.jks" -Force
Copy-Item -Path $truststorePath -Destination "$secretsPath\kafka.truststore.jks" -Force

# Create or update the .env file with the provided values
$envFilePath = ".env"
@"
TRUSTSTORE_LOCATION=/etc/kafka/secrets/kafka.truststore.jks
TRUSTSTORE_PASSWORD=$truststorePasswordPlain
KEYSTORE_LOCATION=/etc/kafka/secrets/kafka.keystore.jks
KEYSTORE_PASSWORD=$keystorePasswordPlain
KEY_PASSWORD=changeit
"@ | Set-Content -Path $envFilePath

Write-Output "Configuration complete. Docker and SSL files have been set up."
