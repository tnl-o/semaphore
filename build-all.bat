@echo off
REM Сборка всего проекта
echo ========================================
echo Сборка Semaphore UI
echo ========================================
echo.

REM Установка Go в PATH
set PATH=C:\Progra~1\Go\bin;%PATH%

REM Сборка Go приложения
echo [1/2] Сборка Go приложения...
cd /d C:\semaphore
go build -o semaphore.exe cli/main.go
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Ошибка сборки Go приложения
    exit /b 1
)
echo   OK: semaphore.exe собран
echo.

REM Сборка фронтенда
echo [2/2] Сборка фронтенда...
cd web
if not exist node_modules (
    echo   Установка зависимостей...
    call npm install
    if %ERRORLEVEL% NEQ 0 (
        echo ERROR: Ошибка установки зависимостей
        exit /b 1
    )
)
echo   Сборка фронтенда...
call npm run build
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Ошибка сборки фронтенда
    exit /b 1
)
echo   OK: Фронтенд собран
cd ..
echo.

echo ========================================
echo Сборка завершена успешно!
echo ========================================
echo.

