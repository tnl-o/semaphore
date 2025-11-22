# PowerShell script to start Semaphore in Docker
# Usage: .\docker-start.ps1

Write-Host "Starting Semaphore in Docker..." -ForegroundColor Green

# Change to project directory
$projectDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectDir

# Start Docker Compose
Write-Host "Building and starting containers..." -ForegroundColor Yellow
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d

if ($LASTEXITCODE -eq 0) {
    Write-Host "`nContainers started successfully!" -ForegroundColor Green
    Write-Host "Waiting for services to be ready..." -ForegroundColor Yellow
    Start-Sleep -Seconds 5
    
    Write-Host "`nChecking container status..." -ForegroundColor Yellow
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
    
    Write-Host "`nViewing logs (last 20 lines)..." -ForegroundColor Yellow
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=20
    
    Write-Host "`n========================================" -ForegroundColor Cyan
    Write-Host "Semaphore should be available at:" -ForegroundColor Green
    Write-Host "http://localhost:3000" -ForegroundColor White
    Write-Host "`nDefault credentials:" -ForegroundColor Green
    Write-Host "Username: admin" -ForegroundColor White
    Write-Host "Password: p455w0rd" -ForegroundColor White
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "`nTo view logs: docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs -f" -ForegroundColor Gray
    Write-Host "To stop: docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down" -ForegroundColor Gray
} else {
    Write-Host "`nFailed to start containers. Check the errors above." -ForegroundColor Red
    Write-Host "`nTrying to view logs..." -ForegroundColor Yellow
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs
}

