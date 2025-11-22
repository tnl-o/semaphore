@echo off
REM Проверка и открытие сервера в браузере
echo ========================================
echo Проверка сервера Semaphore UI
echo ========================================
echo.

REM Проверка, запущен ли сервер
echo Проверка доступности сервера...
timeout /t 3 /nobreak >nul 2>&1

netstat -an | findstr :3000 >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo   OK: Сервер запущен на порту 3000
    echo.
    echo Открываю браузер...
    start http://localhost:3000
    echo.
    echo Сервер доступен по адресу: http://localhost:3000
) else (
    echo   WARNING: Сервер не запущен на порту 3000
    echo.
    echo Для запуска сервера используйте:
    echo   run-dev.bat
    echo   или
    echo   semaphore.exe server --config config.json
)

echo.
pause

