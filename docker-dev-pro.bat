@echo off
chcp 65001 >nul
echo ========================================
echo Запуск Semaphore с PRO функциями для разработки
echo ========================================
echo.

cd /d "%~dp0"

echo [1/5] Остановка существующих контейнеров...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down 2>nul
echo.

echo [2/5] Обновление зависимостей Go...
where go >nul 2>&1
if errorlevel 1 (
    echo Go не установлен локально, обновление зависимостей через Docker...
    docker run --rm -v "%CD%":/workspace -w /workspace golang:1.24-alpine sh -c "apk add --no-cache git && go mod tidy && go mod download"
    if errorlevel 1 (
        echo.
        echo ОШИБКА: Не удалось обновить зависимости через Docker
        echo.
        echo Пожалуйста, сначала запустите:
        echo update-go-deps-docker.bat
        echo.
        pause
        exit /b 1
    )
    echo Зависимости обновлены через Docker
) else (
    go mod tidy
    if errorlevel 1 (
        echo ОШИБКА: Не удалось обновить зависимости
        pause
        exit /b 1
    )
    echo Зависимости обновлены
)

echo.
echo [3/5] Запуск с DEV_MODE=true...
set SEMAPHORE_DEV_MODE=true
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d --build

if errorlevel 1 goto :error

echo.
echo [4/5] Ожидание инициализации базы данных (30 секунд)...
timeout /t 30 /nobreak >nul

echo.
echo [5/5] Активация PRO для пользователя admin...
echo Ожидание готовности базы данных...
set RETRY_COUNT=0
:wait_db
timeout /t 3 /nobreak >nul
docker exec server-server-1 sh -c "test -f /var/lib/semaphore/database.sqlite" >nul 2>&1
if errorlevel 1 (
    set /a RETRY_COUNT+=1
    if %RETRY_COUNT% LSS 10 (
        echo Ожидание базы данных... (%RETRY_COUNT%/10)
        goto wait_db
    ) else (
        echo Предупреждение: База данных не найдена после 30 секунд ожидания
        goto skip_pro
    )
)

echo База данных найдена, обновление пользователя...
docker exec server-server-1 sh -c "sqlite3 /var/lib/semaphore/database.sqlite 'UPDATE user SET pro = 1 WHERE username = \"admin\"; SELECT \"User updated: \" || username || \" (pro=\" || pro || \")\" FROM user WHERE username = \"admin\";'" 2>nul

if errorlevel 1 (
    echo Предупреждение: Не удалось обновить пользователя автоматически
    echo Попробуйте вручную: enable-pro-user.bat
) else (
    echo PRO активирован для пользователя admin
)
:skip_pro

echo.
echo Перезапуск контейнера для применения изменений...
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml restart

echo.
echo ========================================
echo Готово! Semaphore запущен с PRO функциями
echo ========================================
echo.
echo Проверка статуса:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
echo.
echo Последние 20 строк логов:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=20
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
echo ========================================
echo ОШИБКА: Не удалось запустить контейнеры
echo ========================================
echo.
echo Просмотр логов:
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs --tail=50

:end
pause

