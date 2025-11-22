@echo off
chcp 65001 >nul
echo ========================================
echo Исправление проблемы с SQLite в Docker
echo ========================================
echo.

cd /d "%~dp0"

echo [1/3] Остановка контейнеров...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down
if %ERRORLEVEL% NEQ 0 (
    echo Предупреждение: Не удалось остановить контейнеры (возможно, они уже остановлены)
)
echo.

echo [2/3] Удаление старого volume (опционально)...
echo ВНИМАНИЕ: Это удалит все данные из базы данных!
set /p DELETE_VOLUME="Удалить старый volume? (y/n, по умолчанию n): "
if /i "%DELETE_VOLUME%"=="y" (
    echo Удаление volume semaphore_sqlite...
    docker volume rm semaphore_sqlite 2>nul
    if errorlevel 1 (
        echo Volume не найден или уже удален
    ) else (
        echo Volume успешно удален
    )
) else (
    echo Пропущено удаление volume
)
echo.

echo [3/3] Запуск контейнеров с исправленной конфигурацией...
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
    echo Проверка статуса контейнеров:
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
    echo.
    echo Последние 30 строк логов:
    docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=30
    echo.
    echo ========================================
    echo Semaphore должен быть доступен по адресу:
    echo http://localhost:3000
    echo.
    echo Учетные данные:
    echo Username: admin
    echo Password: p455w0rd
    echo ========================================
    echo.
    echo Для просмотра логов в реальном времени:
    echo docker-logs.bat
    goto :end

:error
echo.
echo ========================================
echo ОШИБКА: Не удалось запустить контейнеры
echo ========================================
echo.
echo Просмотр логов для диагностики:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs

:end
pause

