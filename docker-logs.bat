@echo off
REM Batch script to view Semaphore Docker logs
REM Usage: docker-logs.bat

cd /d "%~dp0"

docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs -f

