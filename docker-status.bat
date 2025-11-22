@echo off
REM Проверка статуса Docker контейнера Semaphore
echo ========================================
echo Статус Semaphore Docker контейнера
echo ========================================
echo.

docker ps -a | findstr semaphore
if %ERRORLEVEL% EQU 0 (
    echo.
    echo Логи контейнера:
    echo ----------------------------------------
    docker logs --tail 20 semaphore
) else (
    echo Контейнер не найден
)

echo.
echo Для просмотра логов: docker logs -f semaphore
echo Для остановки: docker-compose down
echo Для запуска: docker-compose up -d
echo.

