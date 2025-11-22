# Настройка Semaphore PRO версии

## Важно

PRO версия Semaphore требует:
1. **GitHub токен** с доступом к приватному репозиторию `semaphoreui/semaphorepro-module`
2. **Сборку из исходников** (нельзя использовать готовый образ)

## Варианты получения PRO версии

### Вариант 1: Использовать готовый PRO образ (если доступен)

Если у вас есть доступ к приватному Docker registry с PRO образами:

```yaml
services:
  server:
    image: docker.io/semaphoreui/semaphore:pro-latest
    # ... остальная конфигурация
```

### Вариант 2: Сборка из исходников с PRO модулем

#### Требования:
- GitHub токен с доступом к `semaphoreui/semaphorepro-module`
- Docker для сборки

#### Шаги:

1. **Получите GitHub токен:**
   - Перейдите в GitHub Settings → Developer settings → Personal access tokens
   - Создайте токен с правами `repo` (для доступа к приватным репозиториям)

2. **Установите переменную окружения:**
   ```cmd
   set GH_TOKEN=your_github_token_here
   ```

3. **Запустите сборку:**
   ```cmd
   docker-start-pro.bat
   ```

   Или вручную:
   ```cmd
   docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build-pro.yml -f deployment/compose/store/sqlite.yml up -d --build
   ```

### Вариант 3: Использовать готовый PRO образ из Docker Hub (если доступен)

Если у вас есть подписка и доступ к PRO образам:

```cmd
docker pull semaphoreui/semaphore:pro-latest
```

Затем измените `base.yml`:
```yaml
image: docker.io/semaphoreui/semaphore:pro-latest
```

## Проверка PRO версии

После запуска проверьте:
1. Откройте http://localhost:3000
2. Войдите в систему
3. PRO функции должны быть доступны (без сообщений "DEMO data")

## Альтернатива: Использовать Community версию

Если у вас нет доступа к PRO модулю, вы можете использовать Community версию, которая имеет большинство функций, но без некоторых PRO-специфичных возможностей.

## Примечания

- PRO версия требует лицензии или подписки
- GitHub токен должен иметь доступ к приватному репозиторию
- Сборка PRO версии занимает больше времени, чем обычная версия
- PRO модуль находится в отдельном приватном репозитории

## Поддержка

Для получения доступа к PRO версии:
- Посетите https://semaphoreui.com
- Свяжитесь с поддержкой для получения лицензии
- Проверьте документацию по лицензированию

