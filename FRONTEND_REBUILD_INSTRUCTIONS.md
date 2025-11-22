# Инструкции по пересборке фронтенда

## Проблема

Новые компоненты созданы и интегрированы в код, но не отображаются в браузере, потому что фронтенд не был пересобран после изменений.

## Что нужно проверить

### 1. Вкладка "Ansible View" отсутствует
- **Ожидается:** В модальном окне задачи должна быть вкладка "Ansible View" между "Log" и "Details"
- **Текущее состояние:** Есть только "Log" и "Details"
- **Причина:** Фронтенд не пересобран

### 2. Breadcrumbs отсутствуют
- **Ожидается:** В верхней части страницы должны быть breadcrumbs (Home → Project → Template)
- **Текущее состояние:** Breadcrumbs не отображаются
- **Причина:** Фронтенд не пересобран

### 3. Ansible Facts не отображаются
- **Ожидается:** Во вкладке "Details" должна быть секция "Ansible Facts" с фактами из логов
- **Текущее состояние:** Нет секции "Ansible Facts"
- **Причина:** Фронтенд не пересобран или facts не извлекаются из логов

## Решение: Пересборка фронтенда

### Вариант 1: Пересборка в Docker (рекомендуется)

Если вы используете Docker, нужно пересобрать образ с обновленным фронтендом:

```bash
# Остановить контейнеры
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml down

# Пересобрать образы (с пересборкой фронтенда)
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml build --no-cache

# Запустить контейнеры
docker-compose -f deployment/compose/server/base.yml -f deployment/compose/server/build.yml -f deployment/compose/store/sqlite.yml up -d
```

### Вариант 2: Локальная пересборка (для разработки)

Если вы разрабатываете локально:

```bash
cd web

# Установить зависимости (если еще не установлены)
npm install

# Пересобрать фронтенд
npm run build

# Или запустить dev-сервер для разработки
npm run serve
```

### Вариант 3: Проверка в dev-режиме

Для быстрой проверки без полной пересборки:

```bash
cd web
npm run serve
```

Это запустит dev-сервер на `http://localhost:8080` (или другом порту), где вы сможете увидеть все изменения в реальном времени.

## После пересборки проверьте

1. **Breadcrumbs:**
   - Откройте любой проект
   - Проверьте верхнюю часть страницы - должны быть breadcrumbs

2. **Вкладка "Ansible View":**
   - Откройте задачу с Ansible playbook
   - В модальном окне должны быть вкладки: "Log", "Ansible View", "Details"
   - Кликните на "Ansible View" - должны быть сгруппированные логи по задачам

3. **Ansible Facts:**
   - Откройте задачу с Ansible playbook
   - Перейдите на вкладку "Details"
   - Прокрутите вниз - должна быть секция "Ansible Facts" (если в логах есть facts)

4. **Консоль браузера:**
   - Откройте DevTools (F12) → Console
   - Не должно быть ошибок загрузки компонентов
   - Проверьте, что нет ошибок типа "Failed to resolve component"

## Проверка версии bundle

После пересборки имя файла bundle изменится (например, с `app.b981070d.js` на `app.xxxxx.js`). Это нормально - это означает, что фронтенд был пересобран.

## Если проблемы остаются

1. **Очистите кэш браузера:**
   - Ctrl+Shift+Delete → Очистить кэш
   - Или используйте режим инкогнито

2. **Проверьте консоль браузера:**
   - F12 → Console
   - Ищите ошибки импорта компонентов

3. **Проверьте Network tab:**
   - F12 → Network
   - Убедитесь, что новый bundle загружается

4. **Проверьте, что компоненты существуют:**
   ```bash
   ls web/src/components/AnsibleLogViewer.vue
   ls web/src/components/Breadcrumbs.vue
   ls web/src/components/AnsibleFactsViewer.vue
   ```

## Быстрая проверка компонентов в коде

Все компоненты должны быть в следующих местах:

- ✅ `web/src/components/AnsibleLogViewer.vue` - существует
- ✅ `web/src/components/Breadcrumbs.vue` - существует  
- ✅ `web/src/components/AnsibleFactsViewer.vue` - должен существовать
- ✅ `web/src/components/TaskLogView.vue` - содержит импорт AnsibleLogViewer
- ✅ `web/src/App.vue` - содержит импорт Breadcrumbs

## Следующие шаги

После пересборки фронтенда все новые компоненты должны появиться в интерфейсе.

