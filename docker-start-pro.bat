@echo off
chcp 65001 >nul
echo ========================================
echo Запуск Semaphore PRO версии
echo ========================================
echo.
echo ВНИМАНИЕ: Для PRO версии требуется:
echo 1. GitHub токен с доступом к semaphoreui/semaphorepro-module
echo 2. Переменная окружения GH_TOKEN
echo.
echo Если у вас нет токена, используйте обычную версию:
echo docker-start.bat
echo.
pause

cd /d "%~dp0"

if "%GH_TOKEN%"=="" (
    echo.
    echo ОШИБКА: Переменная GH_TOKEN не установлена!
    echo.
    echo Установите GitHub токен:
    echo set GH_TOKEN=your_github_token_here
    echo.
    echo Затем запустите скрипт снова.
    pause
    exit /b 1
)

echo [1/2] Остановка существующих контейнеров...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build-pro.yml -f deployment/compose/store/sqlite.yml down 2>nul
echo.

echo [2/2] Сборка и запуск PRO версии...
echo Это может занять несколько минут...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build-pro.yml -f deployment/compose/store/sqlite.yml up -d --build

if errorlevel 1 goto :error

echo.
echo ========================================
echo PRO версия успешно запущена!
echo ========================================
echo.
echo Ожидание инициализации (15 секунд)...
timeout /t 15 /nobreak >nul
echo.
echo Проверка статуса:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build-pro.yml -f deployment/compose/store/sqlite.yml ps
echo.
echo Последние 30 строк логов:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build-pro.yml -f deployment/compose/store/sqlite.yml logs --tail=30
echo.
echo ========================================
echo Откройте в браузере:
echo http://localhost:3000
echo.
echo Учетные данные:
echo Username: admin
echo Password: p455w0rd
echo ========================================
goto :end

:error
echo.
echo ОШИБКА: Не удалось запустить PRO версию
echo.
echo Проверьте:
echo 1. Правильность GitHub токена
echo 2. Доступ к репозиторию semaphoreui/semaphorepro-module
echo 3. Логи ниже:
echo.
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build-pro.yml -f deployment/compose/store/sqlite.yml logs --tail=50

:end
pause

