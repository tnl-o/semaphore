# PowerShell script to fix SQLite issue in Docker
# Usage: .\docker-fix-sqlite.ps1

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Исправление проблемы с SQLite в Docker" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$projectDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $projectDir

Write-Host "[1/3] Остановка контейнеров..." -ForegroundColor Yellow
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down
if ($LASTEXITCODE -ne 0) {
    Write-Host "Предупреждение: Не удалось остановить контейнеры (возможно, они уже остановлены)" -ForegroundColor Yellow
}
Write-Host ""

Write-Host "[2/3] Удаление старого volume (опционально)..." -ForegroundColor Yellow
Write-Host "ВНИМАНИЕ: Это удалит все данные из базы данных!" -ForegroundColor Red
$deleteVolume = Read-Host "Удалить старый volume? (y/n, по умолчанию n)"
if ($deleteVolume -eq "y" -or $deleteVolume -eq "Y") {
    Write-Host "Удаление volume semaphore_sqlite..." -ForegroundColor Yellow
    docker volume rm semaphore_sqlite 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Volume успешно удален" -ForegroundColor Green
    } else {
        Write-Host "Volume не найден или уже удален" -ForegroundColor Yellow
    }
} else {
    Write-Host "Пропущено удаление volume" -ForegroundColor Gray
}
Write-Host ""

Write-Host "[3/3] Запуск контейнеров с исправленной конфигурацией..." -ForegroundColor Yellow
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d

if ($LASTEXITCODE -eq 0) {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "Контейнеры успешно запущены!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Ожидание инициализации (10 секунд)..." -ForegroundColor Yellow
    Start-Sleep -Seconds 10
    Write-Host ""
    Write-Host "Проверка статуса контейнеров:" -ForegroundColor Yellow
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
    Write-Host ""
    Write-Host "Последние 30 строк логов:" -ForegroundColor Yellow
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=30
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host "Semaphore должен быть доступен по адресу:" -ForegroundColor Green
    Write-Host "http://localhost:3000" -ForegroundColor White
    Write-Host ""
    Write-Host "Учетные данные:" -ForegroundColor Green
    Write-Host "Username: admin" -ForegroundColor White
    Write-Host "Password: p455w0rd" -ForegroundColor White
    Write-Host "========================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Для просмотра логов в реальном времени:" -ForegroundColor Gray
    Write-Host "docker-logs.bat" -ForegroundColor Gray
} else {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "ОШИБКА: Не удалось запустить контейнеры" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Просмотр логов для диагностики:" -ForegroundColor Yellow
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs
}

