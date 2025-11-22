@echo off
chcp 65001 >nul
echo ========================================
echo Обновление зависимостей Go через Docker
echo ========================================
echo.

cd /d "%~dp0"

echo Использование Docker для обновления go.mod и go.sum...
echo.

echo Запуск Docker контейнера для обновления зависимостей...
docker run --rm -v "%CD%":/workspace -w /workspace golang:1.24-alpine sh -c "apk add --no-cache git && go mod tidy && go mod download"

if errorlevel 1 (
    echo.
    echo ОШИБКА: Не удалось обновить зависимости через Docker
    echo.
    echo Попробуйте установить Go локально или выполните вручную:
    echo go mod tidy
    echo go mod download
    pause
    exit /b 1
)

echo.
echo ========================================
echo Зависимости успешно обновлены!
echo ========================================
echo.
echo Теперь можно запустить сборку:
echo docker-dev-pro.bat
echo.

pause

