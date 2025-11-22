@echo off
REM Тестирование веб-интерфейса Semaphore
echo ========================================
echo Тестирование веб-интерфейса Semaphore
echo ========================================
echo.

echo [1/3] Проверка контейнера...
docker ps | findstr semaphore
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Контейнер не запущен
    pause
    exit /b 1
)
echo   OK: Контейнер запущен
echo.

echo [2/3] Проверка порта...
netstat -an | findstr :3001 | findstr LISTENING
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Порт 3001 не слушается
    pause
    exit /b 1
)
echo   OK: Порт 3001 доступен
echo.

echo [3/3] Открытие браузера...
start http://localhost:3001
echo   OK: Браузер открыт
echo.

echo ========================================
echo Тестирование завершено
echo ========================================
echo.
echo Веб-интерфейс доступен по адресу: http://localhost:3001
echo.
echo Для просмотра логов: docker logs -f semaphore
echo.

