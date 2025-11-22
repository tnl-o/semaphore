@echo off
chcp 65001 >nul
echo ========================================
echo Активация PRO для пользователя в БД
echo ========================================
echo.

cd /d "%~dp0"

echo Проверка запущенного контейнера...
docker ps | findstr server-server-1 >nul
if errorlevel 1 (
    echo ОШИБКА: Контейнер server-server-1 не запущен!
    echo Сначала запустите: docker-start-dev-pro.bat
    pause
    exit /b 1
)

echo.
echo Обновление пользователя admin в базе данных...
echo.

docker exec server-server-1 sh -c "sqlite3 /var/lib/semaphore/database.sqlite 'UPDATE user SET pro = 1 WHERE username = \"admin\"; SELECT \"User updated: \" || username || \" (pro=\" || pro || \")\" FROM user WHERE username = \"admin\";'"

if errorlevel 1 (
    echo.
    echo ОШИБКА: Не удалось обновить пользователя
    echo Возможно, используется другая БД (не SQLite)
    echo.
    echo Попробуйте вручную через SQL:
    echo docker exec -it server-server-1 sh
    echo sqlite3 /var/lib/semaphore/database.sqlite
    echo UPDATE user SET pro = 1 WHERE username = 'admin';
) else (
    echo.
    echo ========================================
    echo PRO активирован для пользователя admin!
    echo ========================================
    echo.
    echo Перезапустите контейнер для применения изменений:
    echo docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml restart
    echo.
    echo Или просто обновите страницу в браузере
)

pause

