# Инструкция по установке MCP сервера для Rimworld

## Репозиторий
https://github.com/sam2332/mcp-rimworld-dll-search

## Требования

1. **Node.js 14+** и npm
2. **.NET SDK** (для ILSpy CLI)
3. **Установленная игра RimWorld**

## Шаги установки

### 1. Клонирование репозитория

Откройте командную строку (cmd.exe) и выполните:

```cmd
cd C:\Users\abobn\OneDrive\7496~1\semaphore\semaphore
git clone https://github.com/sam2332/mcp-rimworld-dll-search.git
cd mcp-rimworld-dll-search
```

### 2. Установка зависимостей

```cmd
npm install
```

**Примечание:** Если вы видите предупреждения о уязвимостях, можно попробовать исправить их (опционально):
```cmd
npm audit fix
```

Однако это может не исправить все уязвимости, так как некоторые могут быть в зависимостях проекта.

### 3. Установка .NET SDK

**ВАЖНО:** .NET SDK необходим для работы ILSpy CLI. Если вы видите ошибку `"dotnet" не является внутренней или внешней командой`, выполните следующие шаги:

1. **Скачайте .NET SDK:**
   - Перейдите на https://dotnet.microsoft.com/download
   - Скачайте .NET SDK 8.0 или новее (рекомендуется .NET 8.0 SDK)
   - Выберите версию для Windows (x64)

2. **Установите .NET SDK:**
   - Запустите скачанный установщик
   - Следуйте инструкциям установщика
   - Убедитесь, что опция "Add to PATH" включена

3. **Проверьте установку:**
   ```cmd
   dotnet --version
   ```
   Должна отобразиться версия (например, 8.0.xxx)

4. **Установите ILSpy CLI:**
   ```cmd
   dotnet tool install -g ilspycmd
   ```

5. **Проверьте установку ILSpy CLI:**
   ```cmd
   ilspycmd --version
   ```

**Примечание:** После установки .NET SDK может потребоваться перезапустить командную строку, чтобы переменные окружения обновились.

### 4. Настройка пути к RimWorld (опционально)

По умолчанию используется путь: `C:/Program Files (x86)/Steam/steamapps/common/RimWorld`

Если RimWorld установлен в другом месте, установите переменную окружения:

```cmd
set RIMWORLD_PATH=путь\к\вашей\установке\RimWorld
```

### 5. Декомпиляция Assembly-CSharp.dll (первый запуск)

```cmd
npm run decompile
```

Это выполнит:
- Поиск Assembly-CSharp.dll в установке RimWorld
- Декомпиляцию с помощью ILSpy CLI
- Сохранение декомпилированных файлов в папку `./decompile`

### 6. Запуск MCP сервера

```cmd
npm start
```

## Использование MCP сервера

После запуска доступны следующие инструменты:

- **`search`** - Поиск функций
  - Параметры: `searchTerm`, `includePrivate` (опционально), `limit` (опционально)
  
- **`getFullContent`** - Получение полного содержимого файла
  - Параметр: `filePath`

## Интеграция с Cursor

Для использования MCP сервера в Cursor необходимо добавить конфигурацию в настройки Cursor.

Обычно конфигурация MCP серверов находится в файле настроек Cursor. Добавьте следующую конфигурацию:

```json
{
  "mcpServers": {
    "rimworld-dll-search": {
      "command": "node",
      "args": ["C:\\Users\\abobn\\OneDrive\\7496~1\\semaphore\\semaphore\\mcp-rimworld-dll-search\\build\\index.js"],
      "cwd": "C:\\Users\\abobn\\OneDrive\\7496~1\\semaphore\\semaphore\\mcp-rimworld-dll-search"
    }
  }
}
```

**Примечание:** Убедитесь, что путь указан правильно и файл `build/index.js` существует после сборки проекта.

## Дополнительные команды

### Поиск функций через CLI

```cmd
node build/stdioTest.js
```

Команды в интерактивном режиме:
- `search <term> [private=true|false] [limit=number]` - Поиск функций
- `show <filePath>` - Показать полное содержимое файла
- `help` - Показать справку
- `exit` - Выход

### Автоматический тест

```cmd
npm run test
```

Можно изменить условия поиска в файле `src/autoTest.ts`.

## Harmony Patching

После того, как вы найдете функции для патчинга:

1. Получите полное имя класса и сигнатуру функции
2. Создайте Harmony патч используя стандартный синтаксис Harmony
3. Примените патч в инициализации вашего мода

### Пример Harmony патча

```csharp
using HarmonyLib;

[HarmonyPatch(typeof(ThingComp))]
[HarmonyPatch("GetComp")]
[HarmonyPatch(new Type[] { })]
public class Patch_ThingComp_GetComp
{
    public static void Prefix(ThingComp __instance)
    {
        // Ваш код здесь
    }
}
```

## Устранение проблем

### Проблема: ILSpy CLI не найден
**Решение:** Убедитесь, что .NET SDK установлен и ILSpy CLI установлен глобально:
```cmd
dotnet tool install -g ilspycmd
```

### Проблема: Assembly-CSharp.dll не найден
**Решение:** Проверьте путь к RimWorld и установите переменную окружения `RIMWORLD_PATH` если необходимо.

### Проблема: npm install не работает
**Решение:** Убедитесь, что Node.js установлен:
```cmd
node --version
npm --version
```

### Проблема: "dotnet" не является внутренней или внешней командой
**Решение:** 
1. Установите .NET SDK с https://dotnet.microsoft.com/download
2. После установки перезапустите командную строку
3. Проверьте установку: `dotnet --version`
4. Если команда все еще не работает, добавьте путь к .NET SDK вручную в переменную окружения PATH

### Проблема: Error: spawn dotnet ENOENT
**Решение:** Это означает, что .NET SDK не установлен или не добавлен в PATH. См. решение выше.

## Лицензия

MIT

## Благодарности

- RimWorld от Ludeon Studios
- ILSpy для декомпиляции .NET
- Harmony для патчинга .NET во время выполнения

