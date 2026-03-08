@echo off
chcp 65001 >nul
echo ========================================
echo Настройка MCP сервера для Rimworld
echo ========================================
echo.

set "PROJECT_DIR=C:\Users\abobn\OneDrive\7496~1\semaphore\semaphore\mcp-rimworld-dll-search"
set "CURSOR_CONFIG=%APPDATA%\Cursor\User\globalStorage\saoudrizwan.claude-dev\settings\cline_mcp_settings.json"

echo Проверка установки...
echo.

REM Проверка Node.js
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo [ОШИБКА] Node.js не найден! Установите Node.js.
    pause
    exit /b 1
)
echo [OK] Node.js установлен

REM Проверка .NET SDK
where dotnet >nul 2>&1
if %errorlevel% neq 0 (
    echo [ПРЕДУПРЕЖДЕНИЕ] .NET SDK не найден!
    echo Установите .NET SDK с https://dotnet.microsoft.com/download
    echo.
) else (
    echo [OK] .NET SDK установлен
    
    REM Проверка ILSpy CLI
    dotnet tool list -g | findstr ilspycmd >nul 2>&1
    if %errorlevel% neq 0 (
        echo [ИНФО] Установка ILSpy CLI...
        dotnet tool install -g ilspycmd
        if %errorlevel% neq 0 (
            echo [ОШИБКА] Не удалось установить ILSpy CLI
        ) else (
            echo [OK] ILSpy CLI установлен
        )
    ) else (
        echo [OK] ILSpy CLI установлен
    )
)

REM Проверка проекта
if not exist "%PROJECT_DIR%" (
    echo [ОШИБКА] Проект не найден: %PROJECT_DIR%
    echo Выполните: git clone https://github.com/sam2332/mcp-rimworld-dll-search.git
    pause
    exit /b 1
)
echo [OK] Проект найден

REM Переход в директорию проекта
cd /d "%PROJECT_DIR%"

REM Проверка node_modules
if not exist "node_modules" (
    echo [ИНФО] Установка зависимостей...
    call npm install
    if %errorlevel% neq 0 (
        echo [ОШИБКА] Не удалось установить зависимости
        pause
        exit /b 1
    )
    echo [OK] Зависимости установлены
) else (
    echo [OK] Зависимости установлены
)

REM Проверка декомпиляции
if not exist "decompile" (
    echo [ИНФО] Декомпиляция не выполнена. Запускаю декомпиляцию...
    echo Это может занять несколько минут...
    call npm run decompile
    if %errorlevel% neq 0 (
        echo [ОШИБКА] Не удалось выполнить декомпиляцию
        echo Убедитесь, что:
        echo 1. RimWorld установлен
        echo 2. .NET SDK установлен
        echo 3. ILSpy CLI установлен
        pause
        exit /b 1
    )
    echo [OK] Декомпиляция выполнена
) else (
    echo [OK] Декомпиляция уже выполнена
)

REM Проверка build/index.js
if not exist "build\index.js" (
    echo [ИНФО] Сборка проекта...
    call npm run build
    if %errorlevel% neq 0 (
        echo [ОШИБКА] Не удалось собрать проект
        pause
        exit /b 1
    )
    echo [OK] Проект собран
) else (
    echo [OK] Проект собран
)

echo.
echo ========================================
echo Настройка завершена!
echo ========================================
echo.
echo Следующие шаги:
echo 1. Добавьте конфигурацию MCP сервера в Cursor
echo 2. См. файл RIMWORLD_MCP_SETUP.md для инструкций
echo.
echo Конфигурация для Cursor:
echo.
echo {
echo   "mcpServers": {
echo     "rimworld-dll-search": {
echo       "command": "node",
echo       "args": ["%PROJECT_DIR:\=/%\build\index.js"],
echo       "cwd": "%PROJECT_DIR:\=/%",
echo       "env": {
echo         "RIMWORLD_PATH": "C:/Program Files (x86)/Steam/steamapps/common/RimWorld"
echo       }
echo     }
echo   }
echo }
echo.
echo Путь к конфигурации Cursor:
echo %CURSOR_CONFIG%
echo.
pause

