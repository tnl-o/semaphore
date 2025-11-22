@echo off
chcp 65001 >nul
echo ========================================
echo Перезапуск с исправлением WEB_ROOT
echo ========================================
echo.

cd /d "%~dp0"

echo [1/2] Остановка контейнеров...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down
echo.

echo [2/2] Запуск контейнеров с исправленным WEB_ROOT...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d
if errorlevel 1 goto :error

echo.
echo ========================================
echo Контейнеры успешно запущены!
echo ========================================
echo.
echo Ожидание инициализации (10 секунд)...
timeout /t 10 /nobreak >nul
echo.
echo Проверка статуса:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
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
echo ОШИБКА: Не удалось запустить контейнеры
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=30

:end
pause

