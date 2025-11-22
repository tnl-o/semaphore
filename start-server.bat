@echo off
REM Запуск Semaphore сервера
cd /d C:\semaphore
set PATH=C:\Progra~1\Go\bin;%PATH%
set SEMAPHORE_DB_DIALECT=bolt
set SEMAPHORE_PORT=:3000
set SEMAPHORE_TMP_PATH=.\tmp
set SEMAPHORE_WEB_ROOT=http://localhost:3000

if not exist tmp mkdir tmp

echo ========================================
echo Запуск Semaphore Server
echo ========================================
echo.
echo Сервер будет доступен по адресу: http://localhost:3000
echo.
echo Для остановки нажмите Ctrl+C
echo.

semaphore.exe server --no-config

