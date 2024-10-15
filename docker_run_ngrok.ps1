$containerName = "ngrok"
$configFilePath = "S:\Projekti\ngrok\ngrok.yml"

$runningContainer = docker ps -q -f "name=$containerName"
$existingContainer = docker ps -a -q -f "name=$containerName"

if ($runningContainer) {
    Write-Host "Container $containerName is already running."
} elseif ($existingContainer) {
    Write-Host "Starting the existing container: $containerName"
    docker start $containerName
} else {
    Write-Host "No existing container found. Creating and starting a new container: $containerName"
    docker run --name ngrok -p 4040:4040 -p 8081:8081 -v ${configFilePath}:/ngrok.yml ngrok/ngrok http host.docker.internal:8080 --config /ngrok.yml
}