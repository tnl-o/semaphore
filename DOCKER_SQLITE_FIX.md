# Исправление проблемы с SQLite в Docker

## Проблема

При запуске контейнера возникает ошибка:
```
unable to open database file: out of memory (14)
```

Это ошибка SQLite, которая обычно означает проблему с правами доступа к файлу базы данных.

## Причина

Проблема в том, что:
1. Volume монтируется в `/var/lib/database`
2. Но пользователь `semaphore` (UID 1001) может не иметь прав на запись в этот путь
3. В Dockerfile директория `/var/lib/semaphore` уже настроена с правильными правами

## Решение

Изменил конфигурацию SQLite, чтобы использовать `/var/lib/semaphore` вместо `/var/lib/database`, так как эта директория уже настроена с правильными правами в Dockerfile.

## Что было изменено

В файле `deployment/compose/store/sqlite.yml`:
- `SEMAPHORE_DB_PATH` изменен с `/var/lib/database` на `/var/lib/semaphore`
- Volume монтируется в `/var/lib/semaphore`

## Как применить исправление

1. Остановите контейнер:
   ```cmd
   docker-stop.bat
   ```

2. Удалите старый volume (опционально, если хотите начать с чистой базы):
   ```cmd
   docker volume rm semaphore_sqlite
   ```

3. Запустите контейнер заново:
   ```cmd
   docker-start.bat
   ```

## Альтернативное решение

Если проблема сохраняется, можно попробовать использовать PostgreSQL или MySQL вместо SQLite:

### PostgreSQL:
```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/postgres.yml up -d
```

### MySQL:
```cmd
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/mysql.yml up -d
```

