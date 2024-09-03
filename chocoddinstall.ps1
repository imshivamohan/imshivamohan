# Define package paths and filenames
$toolsPath = "$PSScriptRoot\tools"
$dockerComposeFile = "$toolsPath\kafka-ui-docker-compose.yml"
$truststoreFile = "$toolsPath\truststore.jks"
$desktopPath = [System.Environment]::GetFolderPath('Desktop')
$shortcutStartPath = "$desktopPath\Start Kafka UI.lnk"
$shortcutStopPath = "$desktopPath\Stop Kafka UI.lnk"

# Install Rancher Desktop
function Install-RancherDesktop {
    $rancherDesktopUrl = "https://github.com/rancher-sandbox/rancher-desktop/releases/download/v1.0.0/rancher-desktop-1.0.0-windows-amd64.exe"
    $installerPath = "$env:TEMP\rancher-desktop-installer.exe"

    Write-Output "Downloading Rancher Desktop..."
    Invoke-WebRequest -Uri $rancherDesktopUrl -OutFile $installerPath

    Write-Output "Installing Rancher Desktop..."
    Start-Process -FilePath $installerPath -ArgumentList "/SILENT" -Wait

    Remove-Item $installerPath
}

# Validate Docker is running
function Validate-Docker {
    $dockerStatus = docker info 2>&1
    if ($dockerStatus -match "Cannot connect to the Docker daemon") {
        Write-Output "Docker is not running. Please start Docker."
        Exit 1
    }
}

# Copy truststore file
function Copy-Truststore {
    if (Test-Path $truststoreFile) {
        Write-Output "Copying truststore to system..."
        Copy-Item -Path $truststoreFile -Destination "C:\ProgramData\kafka\secrets\truststore.jks" -Force
    } else {
        Write-Output "Truststore file not found."
        Exit 1
    }
}

# Deploy Docker Compose
function Deploy-DockerCompose {
    Write-Output "Deploying Docker Compose setup..."
    docker-compose -f $dockerComposeFile pull
    docker-compose -f $dockerComposeFile up -d
}

# Create shortcut
function Create-Shortcut {
    param (
        [string]$targetPath,
        [string]$shortcutPath,
        [string]$arguments,
        [string]$description
    )
    $WScriptShell = New-Object -ComObject WScript.Shell
    $shortcut = $WScriptShell.CreateShortcut($shortcutPath)
    $shortcut.TargetPath = $targetPath
    $shortcut.Arguments = $arguments
    $shortcut.Description = $description
    $shortcut.Save()
}

# Create shortcuts for start and stop
function Create-Shortcuts {
    $dockerComposeStartCmd = "docker-compose -f `"$dockerComposeFile`" up -d"
    $dockerComposeStopCmd = "docker-compose -f `"$dockerComposeFile`" down"

    Create-Shortcut -targetPath "powershell.exe" `
                    -shortcutPath $shortcutStartPath `
                    -arguments "-NoProfile -ExecutionPolicy Bypass -Command `"$dockerComposeStartCmd`"" `
                    -description "Start Kafka UI"

    Create-Shortcut -targetPath "powershell.exe" `
                    -shortcutPath $shortcutStopPath `
                    -arguments "-NoProfile -ExecutionPolicy Bypass -Command `"$dockerComposeStopCmd`"" `
                    -description "Stop Kafka UI"
}

# Main script execution
Install-RancherDesktop
Validate-Docker
Copy-Truststore
Deploy-DockerCompose
Create-Shortcuts

Write-Output "Kafka UI installation complete. Shortcuts have been created on the desktop."





