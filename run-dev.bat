@echo off
REM Быстрый запуск для разработки (с BoltDB)
echo ========================================
echo Запуск Semaphore UI (Dev Mode)
echo ========================================
echo.

REM Установка Go в PATH
set PATH=C:\Progra~1\Go\bin;%PATH%

REM Проверка билда
if not exist semaphore.exe (
    echo Ошибка: semaphore.exe не найден
    echo Запустите build-all.bat для сборки
    pause
    exit /b 1
)

REM Проверка фронтенда
if not exist web\dist (
    echo Предупреждение: web\dist не найден
    echo Запустите build-all.bat для сборки фронтенда
    pause
)

REM Создание временной директории для данных
if not exist .\tmp (
    mkdir .\tmp
)

REM Запуск с минимальной конфигурацией через переменные окружения
echo Запуск сервера на http://localhost:3000
echo.
echo Для остановки нажмите Ctrl+C
echo.

set SEMAPHORE_DB_DIALECT=bolt
set SEMAPHORE_PORT=:3000
set SEMAPHORE_INTERFACE=
set SEMAPHORE_TMP_PATH=.\tmp
set SEMAPHORE_WEB_ROOT=http://localhost:3000

semaphore.exe server --no-config

