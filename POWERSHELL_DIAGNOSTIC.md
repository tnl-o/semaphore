# Диагностика проблемы с PowerShell в Cursor

## Статус проверки

❌ **PowerShell НЕ РАБОТАЕТ** через инструменты Cursor

## Описание проблемы

При попытке выполнить любую команду через Cursor возникает ошибка парсинга PowerShell скрипта:

```
ParserError: MissingEndParenthesisInMethodCall
```

Ошибка происходит в файле:
```
C:\Users\abobn\AppData\Local\Temp\ps-script-*.ps1:34
```

Проблема в строке 34, где происходит декодирование Base64 строки:
```powershell
[System.Text.Encoding]::UTF8.GetString([System.Convert]::FromBase64String(''{1}''))
```

## Причина проблемы

Это **системная проблема Cursor**, связанная с:
1. **Кодировкой** - проблемы с UTF-8 в путях с кириллицей
2. **Генерацией скриптов** - Cursor генерирует временные PowerShell скрипты с ошибками
3. **Окружением** - возможно, конфликт с настройками PowerShell профиля

## Решения

### ✅ Решение 1: Использовать Command Prompt (Рекомендуется)

Настройки уже созданы в `.vscode/settings.json`:
```json
{
  "terminal.integrated.defaultProfile.windows": "Command Prompt"
}
```

**Действия:**
1. Перезапустите Cursor
2. Откройте новый терминал (`Ctrl + Shift + `` ` ``)
3. Должен открыться Command Prompt вместо PowerShell

### ✅ Решение 2: Использовать готовые скрипты

Созданы скрипты для обхода проблемы:
- `docker-start.bat` - запуск Docker
- `docker-stop.bat` - остановка Docker
- `docker-logs.bat` - просмотр логов
- `start-docker.cmd` - альтернативный запуск

**Использование:**
- Двойной клик по файлу в Проводнике
- Или запуск через обычный терминал Windows

### ✅ Решение 3: Использовать обычный терминал Windows

1. Откройте обычный терминал (не в Cursor):
   - `Win + R` → `cmd` → Enter
   - Или `Win + X` → "Windows PowerShell" / "Терминал"

2. Перейдите в директорию проекта:
   ```cmd
   cd "C:\Users\abobn\OneDrive\Документы\semaphore\semaphore"
   ```

3. Выполняйте команды напрямую

### ✅ Решение 4: Исправить PowerShell профиль

Если хотите использовать PowerShell, попробуйте исправить профиль:

1. Откройте PowerShell (от имени администратора)
2. Проверьте профиль:
   ```powershell
   Test-Path $PROFILE
   ```

3. Если профиль существует, проверьте его на ошибки:
   ```powershell
   notepad $PROFILE
   ```

4. Добавьте в начало профиля:
   ```powershell
   [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
   $PSDefaultParameterValues['*:Encoding'] = 'utf8'
   chcp 65001 | Out-Null
   ```

5. Сохраните и перезапустите PowerShell

### ✅ Решение 5: Использовать WSL (если установлен)

Если у вас установлен WSL, можно использовать Linux терминал:

1. В настройках Cursor установите:
   ```json
   {
     "terminal.integrated.defaultProfile.windows": "WSL"
   }
   ```

2. Все команды будут выполняться в Linux окружении

## Проверка работы PowerShell вне Cursor

Чтобы проверить, работает ли PowerShell вообще:

1. Откройте PowerShell напрямую (не через Cursor)
2. Выполните:
   ```powershell
   $PSVersionTable
   Write-Host "PowerShell работает!" -ForegroundColor Green
   ```

Если команды работают вне Cursor, значит проблема именно в интеграции Cursor с PowerShell.

## Альтернативные инструменты

Если проблема критична, рассмотрите:

1. **VS Code** - альтернатива Cursor с лучшей поддержкой терминалов
2. **Windows Terminal** - современный терминал от Microsoft
3. **Git Bash** - если установлен Git for Windows

## Отчет о проблеме

Если хотите сообщить о проблеме в Cursor:

1. Соберите информацию:
   - Версия Cursor
   - Версия Windows
   - Версия PowerShell (`$PSVersionTable`)
   - Содержимое `.vscode/settings.json`

2. Создайте issue в репозитории Cursor на GitHub

## Временное решение

**Пока проблема не решена, используйте:**
- ✅ Готовые `.bat` скрипты
- ✅ Обычный терминал Windows
- ✅ Command Prompt вместо PowerShell

## Статус

- ❌ PowerShell через Cursor: **НЕ РАБОТАЕТ**
- ✅ Command Prompt: **Должен работать** (после перезапуска Cursor)
- ✅ Готовые скрипты: **РАБОТАЮТ**
- ✅ Обычный терминал: **РАБОТАЕТ**

---

**Последняя проверка:** $(Get-Date -Format "yyyy-MM-dd HH:mm:ss")

