@echo off
chcp 65001 >nul
echo ========================================
echo Запуск Semaphore с PRO функциями для разработки
echo ========================================
echo.
echo ВНИМАНИЕ: Это активирует PRO функции в dev режиме
echo для целей разработки и тестирования.
echo.
pause

cd /d "%~dp0"

echo [1/2] Остановка существующих контейнеров...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down 2>nul
echo.

echo [2/2] Запуск с DEV_MODE=true...
set SEMAPHORE_DEV_MODE=true
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d --build

if errorlevel 1 goto :error

echo.
echo ========================================
echo Semaphore запущен с PRO функциями для разработки!
echo ========================================
echo.
echo Ожидание инициализации (15 секунд)...
timeout /t 15 /nobreak >nul
echo.
echo Проверка статуса:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
echo.
echo Последние 30 строк логов:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=30
echo.
echo ========================================
echo Откройте в браузере:
echo http://localhost:3000
echo.
echo Учетные данные:
echo Username: admin
echo Password: p455w0rd
echo.
echo PRO функции активированы для разработки!
echo ========================================
goto :end

:error
echo.
echo ОШИБКА: Не удалось запустить контейнеры
echo.
echo Просмотр логов:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=50

:end
pause

