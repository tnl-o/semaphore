@echo off
REM Запуск Docker Desktop
echo ========================================
echo Запуск Docker Desktop
echo ========================================
echo.

REM Попытка запустить Docker Desktop
echo Запуск Docker Desktop...
start "" "C:\Program Files\Docker\Docker\Docker Desktop.exe" 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Попытка альтернативного пути...
    start "" "%ProgramFiles%\Docker\Docker\Docker Desktop.exe" 2>nul
)

echo.
echo Ожидание запуска Docker Desktop...
echo Это может занять 30-60 секунд...
echo.

REM Ожидание запуска Docker
:wait_docker
timeout /t 5 /nobreak >nul 2>&1
docker info >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo Docker Desktop запущен!
    echo.
    goto :docker_ready
) else (
    echo Ожидание...
    goto :wait_docker
)

:docker_ready
echo Docker готов к работе.
echo.
