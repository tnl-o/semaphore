@echo off
REM Запуск Semaphore в Docker контейнере
echo ========================================
echo Запуск Semaphore UI в Docker
echo ========================================
echo.

cd /d C:\semaphore

REM Остановка и удаление старого контейнера, если есть
echo [1/3] Остановка старых контейнеров...
docker-compose down 2>nul
echo.

REM Сборка образа (если нужно)
echo [2/3] Сборка Docker образа...
echo Это может занять несколько минут...
docker-compose build
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Ошибка сборки образа
    pause
    exit /b 1
)
echo.

REM Запуск контейнера
echo [3/3] Запуск контейнера...
docker-compose up -d
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Ошибка запуска контейнера
    pause
    exit /b 1
)
echo.

echo ========================================
echo Контейнер запущен!
echo ========================================
echo.
echo Сервер доступен по адресу: http://localhost:3000
echo.
echo Логи: docker-compose logs -f
echo Остановка: docker-compose down
echo.
echo Открываю браузер...
timeout /t 3 /nobreak >nul 2>&1
start http://localhost:3000
echo.

