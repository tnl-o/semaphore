# PRO функции для разработки (Dev Mode)

## Описание

Для целей разработки и тестирования можно активировать PRO функции в Community версии Semaphore, установив переменную окружения `SEMAPHORE_DEV_MODE=true`.

## Что было изменено

1. **`pro/services/server/subscription_svc.go`**:
   - `HasActiveSubscription()` теперь возвращает `true` если `SEMAPHORE_DEV_MODE=true`
   - `CanAddProUser()` теперь возвращает `true` в dev режиме
   - `CanAddRole()`, `CanAddRunner()`, `CanAddTerraformHTTPBackend()` также активированы в dev режиме
   - `GetPlan()` возвращает `"dev"` в dev режиме

2. **`pro/pkg/features/features.go`**:
   - Все PRO функции (`project_runners`, `terraform_backend`, `task_summary`, `secret_storages`) активированы в dev режиме

3. **`deployment/compose/server/base.yml`**:
   - Добавлена переменная окружения `SEMAPHORE_DEV_MODE` (по умолчанию `false`)

## Использование

### Вариант 1: Использовать готовый скрипт

```cmd
docker-start-dev-pro.bat
```

### Вариант 2: Установить переменную окружения вручную

В PowerShell:
```powershell
$env:SEMAPHORE_DEV_MODE="true"
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d --build
```

В Command Prompt:
```cmd
set SEMAPHORE_DEV_MODE=true
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d --build
```

### Вариант 3: Добавить в docker-compose файл

Создайте файл `deployment/compose/server/dev-pro.yml`:
```yaml
services:
  server:
    environment:
      SEMAPHORE_DEV_MODE: "true"
```

Затем запустите:
```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/server/dev-pro.yml -f deployment/compose/store/sqlite.yml up -d --build
```

## Активация PRO для пользователя

После запуска с `SEMAPHORE_DEV_MODE=true` нужно также установить флаг `Pro = true` для пользователя в базе данных.

### Через SQL (для SQLite):

1. Подключитесь к контейнеру:
   ```cmd
   docker exec -it server-server-1 sh
   ```

2. Откройте SQLite базу:
   ```bash
   sqlite3 /var/lib/semaphore/database.sqlite
   ```

3. Обновите пользователя:
   ```sql
   UPDATE user SET pro = 1 WHERE username = 'admin';
   .quit
   ```

### Через API (если есть доступ):

Используйте API endpoint для обновления пользователя, установив `"pro": true` в теле запроса.

## Проверка

После активации:
1. Откройте http://localhost:3000
2. Войдите в систему
3. PRO функции должны быть доступны (без сообщений "DEMO data")
4. Task summary должен работать
5. Другие PRO функции должны быть активны

## Важно

- **Это только для разработки!** Не используйте в production
- PRO функции активированы, но некоторые могут работать в ограниченном режиме
- Это не заменяет реальную PRO версию с полной функциональностью
- Для production используйте официальную PRO версию с лицензией

## Отключение

Чтобы отключить dev режим:
```cmd
set SEMAPHORE_DEV_MODE=false
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml restart
```

Или просто перезапустите без переменной окружения.

