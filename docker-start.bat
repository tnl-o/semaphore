@echo off
REM Batch script to start Semaphore in Docker
REM Usage: docker-start.bat

echo Starting Semaphore in Docker...

cd /d "%~dp0"

echo Building and starting containers...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Containers started successfully!
    echo Waiting for services to be ready...
    timeout /t 5 /nobreak >nul
    
    echo.
    echo Checking container status...
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
    
    echo.
    echo Viewing logs (last 20 lines)...
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=20
    
    echo.
    echo ========================================
    echo Semaphore should be available at:
    echo http://localhost:3000
    echo.
    echo Default credentials:
    echo Username: admin
    echo Password: p455w0rd
    echo ========================================
    echo.
    echo To view logs: docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs -f
    echo To stop: docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down
) else (
    echo.
    echo Failed to start containers. Check the errors above.
    echo.
    echo Trying to view logs...
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs
)

pause

