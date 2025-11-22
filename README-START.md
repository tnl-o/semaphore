# Быстрый запуск Semaphore UI

## Запуск сервера

### Вариант 1: Использование готового скрипта
```cmd
start-server.bat
```

### Вариант 2: Ручной запуск
```cmd
cd C:\semaphore
set SEMAPHORE_DB_DIALECT=bolt
set SEMAPHORE_PORT=:3000
set SEMAPHORE_TMP_PATH=.\tmp
set SEMAPHORE_WEB_ROOT=http://localhost:3000
semaphore.exe server --no-config
```

### Вариант 3: С конфигурацией
```cmd
semaphore.exe setup
semaphore.exe server --config config.json
```

## Проверка работы

1. Откройте браузер и перейдите на: http://localhost:3000
2. Проверьте, что сервер отвечает на запросы

## Остановка сервера

Нажмите `Ctrl+C` в окне, где запущен сервер.

## Проверка статуса

Используйте скрипт `test-server.bat` для проверки, запущен ли сервер:
```cmd
test-server.bat
```

Или проверьте порт вручную:
```cmd
netstat -an | findstr :3000
```

