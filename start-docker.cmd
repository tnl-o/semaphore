@echo off
chcp 65001 >nul
cd /d "%~dp0"
echo Starting Semaphore Docker containers...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d
if %ERRORLEVEL% EQU 0 (
    echo.
    echo Containers started! Waiting 5 seconds...
    timeout /t 5 /nobreak >nul
    echo.
    echo Container status:
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
    echo.
    echo Last 20 lines of logs:
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=20
    echo.
    echo ========================================
    echo Semaphore is available at: http://localhost:3000
    echo Username: admin
    echo Password: p455w0rd
    echo ========================================
) else (
    echo.
    echo Failed to start containers. Showing logs:
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs
)
pause


