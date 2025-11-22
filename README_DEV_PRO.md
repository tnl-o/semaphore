# Запуск Semaphore с PRO функциями для разработки

## Быстрый старт

### Вариант 1: Автоматический (рекомендуется)

1. **Обновите зависимости Go:**
   ```cmd
   .\update-go-deps-docker.bat
   ```

2. **Запустите Semaphore с PRO функциями:**
   ```cmd
   .\docker-dev-pro.bat
   ```

### Вариант 2: Если Go установлен локально

1. **Обновите зависимости:**
   ```cmd
   go mod tidy
   go mod download
   ```

2. **Запустите:**
   ```cmd
   .\docker-dev-pro.bat
   ```

## Что делает скрипт

`docker-dev-pro.bat` автоматически:
1. Останавливает существующие контейнеры
2. Обновляет зависимости Go (через Docker, если Go не установлен)
3. Запускает Semaphore с `SEMAPHORE_DEV_MODE=true`
4. Ждет инициализации базы данных
5. Активирует PRO для пользователя `admin`
6. Перезапускает контейнер

## Проверка

После запуска:
- Откройте http://localhost:3000
- Войдите как `admin` / `p455w0rd`
- PRO функции должны быть доступны

## Устранение проблем

### Ошибка: "missing go.sum entry"

Если видите ошибку про `go.sum`, выполните:
```cmd
.\update-go-deps-docker.bat
```

Затем запустите снова:
```cmd
.\docker-dev-pro.bat
```

### Ошибка при сборке Docker

Убедитесь, что:
- Docker запущен
- Достаточно места на диске
- Порты 3000 не заняты

### PRO функции не работают

1. Проверьте, что `SEMAPHORE_DEV_MODE=true` установлен:
   ```cmd
   docker exec server-server-1 env | findstr DEV_MODE
   ```

2. Проверьте, что пользователь имеет `pro = true`:
   ```cmd
   docker exec server-server-1 sqlite3 /var/lib/semaphore/database.sqlite "SELECT username, pro FROM user WHERE username = 'admin';"
   ```

3. Если `pro = 0`, выполните:
   ```cmd
   .\enable-pro-user.bat
   ```

## Отключение PRO режима

Чтобы отключить dev режим:
```cmd
set SEMAPHORE_DEV_MODE=false
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml restart
```

Или просто перезапустите без переменной окружения.

