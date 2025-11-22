@echo off
REM Batch script to stop Semaphore Docker containers
REM Usage: docker-stop.bat

echo Stopping Semaphore containers...

cd /d "%~dp0"

docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down

if %ERRORLEVEL% EQU 0 (
    echo Containers stopped successfully!
) else (
    echo Failed to stop containers.
)

pause

