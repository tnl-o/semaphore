# Запуск Semaphore в Docker

## Быстрый старт

### Windows (Batch файлы)

1. **Запуск:**
   ```cmd
   docker-start.bat
   ```

2. **Просмотр логов:**
   ```cmd
   docker-logs.bat
   ```

3. **Остановка:**
   ```cmd
   docker-stop.bat
   ```

### Windows (PowerShell)

1. **Запуск:**
   ```powershell
   .\docker-start.ps1
   ```

### Linux/Mac (Bash)

1. **Запуск:**
   ```bash
   docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d
   ```

2. **Просмотр логов:**
   ```bash
   docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs -f
   ```

3. **Остановка:**
   ```bash
   docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down
   ```

## Доступ к приложению

После запуска приложение будет доступно по адресу:
- **URL:** http://localhost:3000

**Учетные данные по умолчанию:**
- Username: `admin`
- Password: `p455w0rd`
- Email: `admin@localhost`

## Изменение учетных данных

Вы можете изменить учетные данные через переменные окружения перед запуском:

```cmd
set SEMAPHORE_ADMIN_USERNAME=myadmin
set SEMAPHORE_ADMIN_PASSWORD=mypassword
set SEMAPHORE_ADMIN_EMAIL=admin@example.com
docker-start.bat
```

## Проверка статуса

Проверить статус контейнеров:
```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml ps
```

## Использование готового образа (без сборки)

Если сборка занимает слишком много времени, можно использовать готовый образ:

```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/store/sqlite.yml up -d
```

## Использование PostgreSQL вместо SQLite

Для использования PostgreSQL:

```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/postgres.yml up -d
```

## Устранение неполадок

### Порт 3000 занят

Если порт 3000 уже занят, измените порт в файле `deployment/compose/server/base.yml`:

```yaml
ports:
  - "3001:3000"  # Измените 3000 на 3001
```

### Проблемы со сборкой

Если сборка не удается, попробуйте:
1. Очистить Docker кэш: `docker system prune -a`
2. Использовать готовый образ (см. выше)
3. Проверить логи: `docker-logs.bat`

### Контейнер не запускается

Проверьте логи:
```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml logs
```

### Очистка

Удалить все контейнеры и volumes:
```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down -v
```

