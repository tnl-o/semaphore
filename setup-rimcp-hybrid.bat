@echo off
setlocal enabledelayedexpansion
chcp 65001 >nul
echo ========================================
echo Установка RiMCP_hybrid для Rimworld
echo ========================================
echo.

set "PROJECT_DIR=C:\Users\abobn\OneDrive\7496~1\semaphore\semaphore\RiMCP_hybrid"
set "RIMWORLD_PATH=C:\Program Files (x86)\Steam\steamapps\common\RimWorld"
set "RIMWORLD_DATA=!PROJECT_DIR!\RimWorldData"

echo [1/8] Проверка .NET SDK...
where dotnet >nul 2>&1
if %errorlevel% neq 0 (
    echo   [ОШИБКА] .NET SDK не найден!
    echo   Установите .NET 8.0 SDK с https://dotnet.microsoft.com/download
    pause
    exit /b 1
)
echo   [OK] .NET SDK установлен
dotnet --version
echo.

echo [2/8] Проверка Python...
where python >nul 2>&1
if %errorlevel% neq 0 (
    echo   [ОШИБКА] Python не найден!
    echo   Установите Python 3.9+ с https://www.python.org/
    pause
    exit /b 1
)
echo   [OK] Python установлен
python --version
echo.

echo [3/8] Клонирование репозитория...
if exist "%PROJECT_DIR%" (
    echo   [ИНФО] Проект уже существует, пропускаю клонирование
) else (
    cd /d "C:\Users\abobn\OneDrive\7496~1\semaphore\semaphore"
    git clone https://github.com/h7lu/RiMCP_hybrid.git
    if %errorlevel% neq 0 (
        echo   [ОШИБКА] Не удалось клонировать репозиторий
        pause
        exit /b 1
    )
    echo   [OK] Репозиторий клонирован
)
echo.

echo [4/8] Проверка структуры проекта...
cd /d "!PROJECT_DIR!"
if not exist "src" (
    echo   [ОШИБКА] Структура проекта неверна!
    pause
    exit /b 1
)
echo   [OK] Структура проекта корректна
echo.

echo [5/8] Проверка данных RimWorld...
if not exist "!RIMWORLD_DATA!" (
    echo   [ИНФО] Создание папки RimWorldData...
    mkdir "!RIMWORLD_DATA!"
)

REM Проверка RimWorld используя dir для путей с пробелами
set "RIMWORLD_EXE=!RIMWORLD_PATH!\RimWorld.exe"
dir /b "!RIMWORLD_EXE!" >nul 2>&1
if !errorlevel! equ 0 (
    echo   [OK] RimWorld найден
    if not exist "!RIMWORLD_DATA!\Data" (
        echo   [ИНФО] Копирование Def данных...
        xcopy "!RIMWORLD_PATH!\Data" "!RIMWORLD_DATA!\Data\" /E /I /Y >nul 2>&1
        if !errorlevel! equ 0 (
            echo   [OK] Def данные скопированы
        ) else (
            echo   [ПРЕДУПРЕЖДЕНИЕ] Не удалось скопировать Def данные автоматически
            echo   Скопируйте вручную: "!RIMWORLD_PATH!\Data" -^> "!RIMWORLD_DATA!\Data"
        )
    ) else (
        echo   [OK] Def данные уже существуют
    )
) else (
    echo   [ПРЕДУПРЕЖДЕНИЕ] RimWorld не найден по стандартному пути
    echo   Путь: !RIMWORLD_PATH!
    echo   Вам нужно будет вручную скопировать данные
)
echo.

echo [6/8] Настройка Python окружения...
if exist "scripts\setup-embedding-env.ps1" (
    echo   [ИНФО] Запуск setup-embedding-env.ps1...
    powershell -ExecutionPolicy Bypass -File "scripts\setup-embedding-env.ps1"
    if %errorlevel% neq 0 (
        echo   [ПРЕДУПРЕЖДЕНИЕ] Возможны проблемы с настройкой Python окружения
        echo   Попробуйте запустить вручную: .\scripts\setup-embedding-env.ps1
    ) else (
        echo   [OK] Python окружение настроено
    )
) else (
    echo   [ПРЕДУПРЕЖДЕНИЕ] Скрипт setup-embedding-env.ps1 не найден
    echo   Проверьте структуру проекта
)
echo.

echo [7/8] Сборка проекта...
dotnet build
if %errorlevel% neq 0 (
    echo   [ОШИБКА] Не удалось собрать проект
    pause
    exit /b 1
)
echo   [OK] Проект собран
echo.

echo [8/8] Проверка готовности к созданию индекса...
if not exist "!RIMWORLD_DATA!\Data" (
    echo   [ПРЕДУПРЕЖДЕНИЕ] Def данные не найдены!
    echo   Необходимо скопировать:
    echo   "!RIMWORLD_PATH!\Data" -^> "!RIMWORLD_DATA!\Data"
    echo.
)

echo ========================================
echo Установка завершена!
echo ========================================
echo.
echo Следующие шаги:
echo.
echo 1. Убедитесь, что RimWorldData содержит:
echo    - Data\ (Def файлы из RimWorld)
echo    - C# исходники (декомпилированные .cs файлы)
echo.
echo 2. Если у вас уже есть декомпилированные файлы из mcp-rimworld-dll-search:
echo    Скопируйте их в RimWorldData\
echo.
echo 3. Создайте индекс:
echo    cd src\RimWorldCodeRag
echo    dotnet run -- index --root "..\..\RimWorldData"
echo.
echo 4. Запустите embedding сервер (терминал 1):
echo    .\scripts\start-embedding-server.ps1
echo.
echo 5. Запустите MCP сервер (терминал 2):
echo    cd src\RimWorldCodeRag.McpServer
echo    dotnet run
echo.
echo 6. Добавьте конфигурацию в Cursor (см. RIMCP_HYBRID_INSTALL.md)
echo.
pause

