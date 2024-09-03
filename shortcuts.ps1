# Define paths and project details
$projectPath = "C:\Project" # Adjust this to your project path
$truststoreSourcePath = "$projectPath\truststore\truststore.jks"
$dockerComposeFile = "$projectPath\kafka-ui-docker-compose.yml"
$envFilePath = "$projectPath\.env"
$secretsPath = "C:\ProgramData\kafka\secrets"
$desktopPath = [System.Environment]::GetFolderPath("Desktop")
$startScriptPath = "$projectPath\start-kafka-ui.ps1"
$stopScriptPath = "$projectPath\stop-kafka-ui.ps1"

# (Previous steps remain the same...)

# Create the start script
@"
Start-Process -NoNewWindow -FilePath "docker-compose" -ArgumentList "-f `"$dockerComposeFile`" up -d" -WorkingDirectory "$projectPath"
"@ | Set-Content -Path $startScriptPath

# Create the stop script
@"
Start-Process -NoNewWindow -FilePath "docker-compose" -ArgumentList "-f `"$dockerComposeFile`" down" -WorkingDirectory "$projectPath"
"@ | Set-Content -Path $stopScriptPath

# Create Start Shortcut on Desktop
$startShortcutPath = "$desktopPath\Start Kafka UI.lnk"
$startWScript = New-Object -ComObject WScript.Shell
$startShortcut = $startWScript.CreateShortcut($startShortcutPath)
$startShortcut.TargetPath = "powershell.exe"
$startShortcut.Arguments = "-ExecutionPolicy Bypass -File `"$startScriptPath`""
$startShortcut.WorkingDirectory = $projectPath
$startShortcut.IconLocation = "powershell.exe"
$startShortcut.Save()

# Create Stop Shortcut on Desktop
$stopShortcutPath = "$desktopPath\Stop Kafka UI.lnk"
$stopWScript = New-Object -ComObject WScript.Shell
$stopShortcut = $stopWScript.CreateShortcut($stopShortcutPath)
$stopShortcut.TargetPath = "powershell.exe"
$stopShortcut.Arguments = "-ExecutionPolicy Bypass -File `"$stopScriptPath`""
$stopShortcut.WorkingDirectory = $projectPath
$stopShortcut.IconLocation = "powershell.exe"
$stopShortcut.Save()

Write-Output "Kafka UI has been installed. Shortcuts to start and stop the service have been created on the desktop."
