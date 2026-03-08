# Установка RiMCP_hybrid - Продвинутый MCP сервер для Rimworld

## Описание

[RiMCP_hybrid](https://github.com/h7lu/RiMCP_hybrid) - это более продвинутый MCP сервер для Rimworld, который использует:
- **RAG (Retrieval-Augmented Generation)** с семантическим поиском
- **Embedding модели** для понимания контекста
- **Граф зависимостей** для анализа кода
- **Смешанный поиск** (семантический + точный)

## Преимущества перед простым поиском

- ✅ Понимает естественный язык запросов
- ✅ Семантический поиск (поиск по смыслу, а не только по ключевым словам)
- ✅ Анализ зависимостей (что использует символ, что использует символ)
- ✅ Работает с XML определениями и C# кодом
- ✅ Более точные результаты

## Требования

1. **.NET 8.0 SDK** (уже установлен ✅)
2. **Python 3.9+** (нужно проверить)
3. **RimWorld установлен**
4. **Данные RimWorld** (Def файлы и C# исходники)

## Шаги установки

### 1. Клонирование репозитория

```cmd
cd C:\Users\abobn\OneDrive\7496~1\semaphore\semaphore
git clone https://github.com/h7lu/RiMCP_hybrid.git
cd RiMCP_hybrid
```

### 2. Подготовка данных RimWorld

#### 2.1. Копирование Def данных

Скопируйте папку `Data` из установки RimWorld:
- **Источник:** `C:\Program Files (x86)\Steam\steamapps\common\RimWorld\Data`
- **Назначение:** `RiMCP_hybrid\RimWorldData\Data`

#### 2.2. Экспорт C# исходников

Используйте ILSpy или dnSpy для экспорта исходников:
- **Источник:** `Assembly-CSharp.dll` из RimWorld
- **Назначение:** `RiMCP_hybrid\RimWorldData\` (папка с .cs файлами)

**Примечание:** Если у вас уже есть декомпилированные файлы из предыдущего проекта, можно использовать их.

### 3. Настройка Python окружения (одноразово)

```cmd
cd RiMCP_hybrid
.\scripts\setup-embedding-env.ps1
```

Этот скрипт:
- Создаст виртуальное окружение Python
- Установит зависимости
- Скачает embedding модель (e5-base-v2)

### 4. Сборка проекта

```cmd
dotnet build
```

### 5. Создание индекса (одноразово)

```cmd
cd src\RimWorldCodeRag
dotnet run -- index --root "..\..\RimWorldData"
```

Это создаст индекс для поиска. Может занять несколько минут.

### 6. Запуск сервера

#### Терминал 1: Embedding сервер

```cmd
.\scripts\start-embedding-server.ps1
```

Дождитесь сообщения "Model loaded successfully"

#### Терминал 2: MCP сервер

```cmd
cd src\RimWorldCodeRag.McpServer
dotnet run
```

## Настройка в Cursor

### Вариант 1: Через переменные окружения

Создайте или отредактируйте конфигурацию MCP в Cursor:

```json
{
  "mcpServers": {
    "rimworld-code-rag": {
      "command": "dotnet",
      "args": [
        "run",
        "--project",
        "C:\\Users\\abobn\\OneDrive\\7496~1\\semaphore\\semaphore\\RiMCP_hybrid\\src\\RimWorldCodeRag.McpServer"
      ],
      "env": {
        "RIMWORLD_INDEX_ROOT": "C:\\Users\\abobn\\OneDrive\\7496~1\\semaphore\\semaphore\\RiMCP_hybrid\\index",
        "EMBEDDING_SERVER_URL": "http://127.0.0.1:5000"
      }
    }
  }
}
```

### Вариант 2: Через appsettings.json

Отредактируйте `src/RimWorldCodeRag.McpServer/appsettings.json`:

```json
{
  "McpServer": {
    "IndexRoot": "C:/Users/abobn/OneDrive/7496~1/semaphore/semaphore/RiMCP_hybrid/index",
    "EmbeddingServerUrl": "http://127.0.0.1:5000"
  }
}
```

## Доступные инструменты MCP

### 1. `rough_search` - Семантический поиск
Использует естественный язык для поиска символов и XML определений.

**Пример:**
```
Найди все функции, связанные с системой голода колонистов
```

### 2. `get_uses` - Анализ зависимостей (вниз по потоку)
Показывает, что использует символ.

**Пример:**
```
Покажи, что использует класс Pawn
```

### 3. `get_used_by` - Обратный анализ зависимостей (вверх по потоку)
Показывает, что использует символ.

**Пример:**
```
Покажи, что вызывает метод TakeDamage
```

### 4. `get_item` - Получение исходного кода
Получает полный исходный код символа.

**Пример:**
```
Покажи полный код класса Verse.Pawn
```

## Устранение проблем

### "Index not found"
- Убедитесь, что выполнили шаг 5 (создание индекса)
- Проверьте путь в `RIMWORLD_INDEX_ROOT`

### "Embedding server connection failed"
- Запустите embedding сервер в отдельном терминале
- Дождитесь "Model loaded successfully"
- Проверьте, что порт 5000 свободен

### Python не найден
- Установите Python 3.9+ с https://www.python.org/
- Убедитесь, что Python добавлен в PATH

### Ошибки при создании индекса
- Проверьте, что `RimWorldData` содержит файлы
- Убедитесь, что .NET 8.0 SDK установлен

## Производительность

- **Холодный старт:** ~2-5 секунд (загрузка индекса)
- **Горячий запрос:** 0.5-1 секунда
- **Использование памяти:** ~300MB для векторного индекса
- **Рекомендуется GPU** для embedding сервера (значительное ускорение)

## Обновление после обновления RimWorld

1. Обновите файлы в `RimWorldData/`
2. Пересоздайте индекс:
   ```cmd
   dotnet run -- index --root "..\..\RimWorldData" --force
   ```

## Сравнение с простым поиском

| Функция | mcp-rimworld-dll-search | RiMCP_hybrid |
|---------|------------------------|--------------|
| Простой поиск | ✅ | ✅ |
| Семантический поиск | ❌ | ✅ |
| Понимание естественного языка | ❌ | ✅ |
| Анализ зависимостей | ❌ | ✅ |
| Работа с XML Def | ❌ | ✅ |
| Требует Python | ❌ | ✅ |
| Требует embedding модель | ❌ | ✅ |
| Сложность установки | Простая | Средняя |

## Рекомендации

- Если нужен **простой поиск по коду** - используйте `mcp-rimworld-dll-search`
- Если нужен **продвинутый семантический поиск и анализ** - используйте `RiMCP_hybrid`

Оба сервера можно использовать одновременно, настроив их в Cursor с разными именами.

