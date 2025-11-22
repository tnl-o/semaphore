@echo off
REM Быстрая проверка готовности проекта к запуску
echo ========================================
echo Проверка готовности Semaphore UI
echo ========================================
echo.

REM Проверка Go билда
echo [1/4] Проверка Go билда...
if exist semaphore.exe (
    echo   OK: semaphore.exe найден
) else (
    echo   WARNING: semaphore.exe не найден, нужно собрать: go build -o semaphore.exe cli/main.go
)
echo.

REM Проверка фронтенда
echo [2/4] Проверка фронтенда...
set FRONTEND_FOUND=0
if exist api\public (
    echo   OK: api\public найден (фронтенд собран)
    set FRONTEND_FOUND=1
)
if %FRONTEND_FOUND%==0 if exist web\dist (
    echo   OK: web\dist найден
    set FRONTEND_FOUND=1
)
if %FRONTEND_FOUND%==0 (
    echo   WARNING: Фронтенд не собран, нужно собрать: cd web ^&^& npm install ^&^& npm run build
)
echo.

REM Проверка Go тестов
echo [3/4] Проверка Go тестов...
set PATH=C:\Progra~1\Go\bin;%PATH%
go test ./... >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   OK: Все тесты проходят
) else (
    echo   ERROR: Тесты не проходят, запустите: go test ./...
)
echo.

REM Проверка зависимостей фронтенда
echo [4/4] Проверка зависимостей фронтенда...
if exist web\node_modules (
    echo   OK: node_modules найден
) else (
    echo   WARNING: node_modules не найден, нужно установить: cd web ^&^& npm install
)
echo.

echo ========================================
echo Проверка завершена
echo ========================================
echo.
echo Для запуска сервиса:
echo   1. Настройте конфигурацию: semaphore.exe setup
echo   2. Запустите сервер: semaphore.exe server --config config.json
echo.

