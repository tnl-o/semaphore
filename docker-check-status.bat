@echo off
chcp 65001 >nul
echo ========================================
echo Проверка статуса Docker контейнеров
echo ========================================
echo.

cd /d "%~dp0"

echo [1] Статус контейнеров:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
echo.

echo [2] Последние 50 строк логов:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=50
echo.

echo [3] Проверка доступности порта 3000:
netstat -ano | findstr :3000
echo.

echo [4] Проверка контейнера изнутри (если запущен):
docker exec server-server-1 ps aux 2>nul || echo Контейнер не запущен или недоступен
echo.

pause

