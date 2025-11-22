@echo off
chcp 65001 >nul
echo ========================================
echo Обновление зависимостей Go
echo ========================================
echo.

cd /d "%~dp0"

echo Обновление go.mod и go.sum...
go mod tidy

if errorlevel 1 (
    echo.
    echo ОШИБКА: Не удалось обновить зависимости
    echo Убедитесь, что Go установлен и доступен в PATH
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

