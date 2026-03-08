@echo off
chcp 65001 >nul
echo ========================================
echo Проверка подключения MCP сервера Rimworld
echo ========================================
echo.

set "PROJECT_DIR=C:\Users\abobn\OneDrive\7496~1\semaphore\semaphore\mcp-rimworld-dll-search"
set "RIMWORLD_PATH=C:/Program Files (x86)/Steam/steamapps/common/RimWorld"

echo [1/6] Проверка существования проекта...
if exist "%PROJECT_DIR%" (
    echo   [OK] Проект найден: %PROJECT_DIR%
) else (
    echo   [ОШИБКА] Проект не найден!
    echo   Ожидаемый путь: %PROJECT_DIR%
    pause
    exit /b 1
)
echo.

echo [2/6] Проверка структуры проекта...
cd /d "%PROJECT_DIR%"

if exist "package.json" (
    echo   [OK] package.json найден
) else (
    echo   [ОШИБКА] package.json не найден!
    pause
    exit /b 1
)

if exist "build\index.js" (
    echo   [OK] build\index.js найден
) else (
    echo   [ПРЕДУПРЕЖДЕНИЕ] build\index.js не найден
    echo   Попытка сборки...
    call npm run build
    if %errorlevel% neq 0 (
        echo   [ОШИБКА] Не удалось собрать проект
        pause
        exit /b 1
    )
)
echo.

echo [3/6] Проверка декомпилированных файлов...
if exist "decompile" (
    echo   [OK] Папка decompile существует
    dir /b "decompile" 2>nul | find /c /v "" > temp_count.txt
    set /p FILE_COUNT=<temp_count.txt
    del temp_count.txt
    if !FILE_COUNT! GTR 0 (
        echo   [OK] Найдено файлов в decompile: !FILE_COUNT!
    ) else (
        echo   [ПРЕДУПРЕЖДЕНИЕ] Папка decompile пуста
        echo   Необходимо выполнить декомпиляцию
    )
) else (
    echo   [ПРЕДУПРЕЖДЕНИЕ] Папка decompile не существует
    echo   Необходимо выполнить декомпиляцию
)
echo.

echo [4/6] Проверка пути к RimWorld...
if exist "%RIMWORLD_PATH%\RimWorld.exe" (
    echo   [OK] RimWorld найден: %RIMWORLD_PATH%
) else (
    echo   [ПРЕДУПРЕЖДЕНИЕ] RimWorld не найден по стандартному пути
    echo   Проверьте путь: %RIMWORLD_PATH%
    echo   Если RimWorld установлен в другом месте, установите переменную:
    echo   set RIMWORLD_PATH=ваш_путь
)
echo.

echo [5/6] Проверка .NET SDK и ILSpy CLI...
where dotnet >nul 2>&1
if %errorlevel% neq 0 (
    echo   [ОШИБКА] .NET SDK не найден!
) else (
    echo   [OK] .NET SDK установлен
    dotnet --version
    dotnet tool list -g | findstr ilspycmd >nul 2>&1
    if %errorlevel% neq 0 (
        echo   [ПРЕДУПРЕЖДЕНИЕ] ILSpy CLI не установлен
        echo   Установка ILSpy CLI...
        dotnet tool install -g ilspycmd
    ) else (
        echo   [OK] ILSpy CLI установлен
    )
)
echo.

echo [6/6] Проверка конфигурации Cursor...
set "CURSOR_CONFIG=%APPDATA%\Cursor\User\globalStorage\saoudrizwan.claude-dev\settings\cline_mcp_settings.json"
if exist "%CURSOR_CONFIG%" (
    echo   [OK] Конфигурация Cursor найдена
    echo   Путь: %CURSOR_CONFIG%
    echo.
    echo   Содержимое конфигурации:
    type "%CURSOR_CONFIG%"
) else (
    echo   [ПРЕДУПРЕЖДЕНИЕ] Конфигурация Cursor не найдена
    echo   Ожидаемый путь: %CURSOR_CONFIG%
)
echo.

echo ========================================
echo Результаты проверки
echo ========================================
echo.
echo Если декомпиляция не выполнена, выполните:
echo   cd "%PROJECT_DIR%"
echo   npm run decompile
echo.
echo Если декомпиляция не работает, проверьте:
echo   1. Путь к RimWorld правильный
echo   2. .NET SDK установлен
echo   3. ILSpy CLI установлен
echo   4. Assembly-CSharp.dll существует в RimWorld
echo.
pause

