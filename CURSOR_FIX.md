# Исправление окружения Cursor для Windows

## Проблема

Cursor использует PowerShell для выполнения команд, но возникают ошибки парсинга из-за проблем с кодировкой или конфигурацией PowerShell.

## Решения

### Решение 1: Использовать Command Prompt вместо PowerShell (Рекомендуется)

1. Откройте настройки Cursor:
   - `Ctrl + ,` (или `File > Preferences > Settings`)
   - Или откройте файл `.vscode/settings.json`

2. Убедитесь, что в файле `.vscode/settings.json` установлено:
   ```json
   {
     "terminal.integrated.defaultProfile.windows": "Command Prompt"
   }
   ```

3. Перезапустите Cursor

### Решение 2: Исправить настройки PowerShell

Если вы хотите использовать PowerShell, выполните в PowerShell (от имени администратора):

```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
[Console]::OutputEncoding = [System.Text.Encoding]::UTF8
$PSDefaultParameterValues['*:Encoding'] = 'utf8'
```

### Решение 3: Использовать Git Bash

1. Установите Git for Windows (если еще не установлен)
2. В настройках Cursor установите:
   ```json
   {
     "terminal.integrated.defaultProfile.windows": "Git Bash"
   }
   ```

### Решение 4: Исправить кодировку PowerShell

Создайте файл профиля PowerShell:

1. Откройте PowerShell
2. Выполните:
   ```powershell
   if (!(Test-Path -Path $PROFILE)) {
       New-Item -ItemType File -Path $PROFILE -Force
   }
   ```

3. Добавьте в профиль (`notepad $PROFILE`):
   ```powershell
   [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
   $PSDefaultParameterValues['*:Encoding'] = 'utf8'
   chcp 65001 | Out-Null
   ```

4. Перезапустите PowerShell

### Решение 5: Использовать WSL (Windows Subsystem for Linux)

Если у вас установлен WSL:

1. В настройках Cursor установите:
   ```json
   {
     "terminal.integrated.defaultProfile.windows": "WSL"
   }
   ```

2. Все команды будут выполняться в Linux окружении

## Проверка

После применения изменений:

1. Перезапустите Cursor
2. Откройте новый терминал (`Ctrl + Shift + `` ` ``)
3. Проверьте, что терминал работает:
   ```cmd
   echo Test
   ```

## Альтернатива: Использовать готовые скрипты

Если проблемы с терминалом продолжаются, используйте созданные скрипты:

- `docker-start.bat` - для запуска Docker
- `docker-stop.bat` - для остановки Docker
- `docker-logs.bat` - для просмотра логов

Запускайте их напрямую в обычном терминале Windows (не в Cursor).

## Дополнительные настройки

### Настройка кодировки для проекта

Создайте файл `.editorconfig` в корне проекта (если его нет):

```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true
```

### Настройка Git

Убедитесь, что Git настроен правильно:

```bash
git config --global core.autocrlf true
git config --global core.quotepath false
```

## Если ничего не помогает

1. Обновите Cursor до последней версии
2. Переустановите Cursor
3. Используйте обычный терминал Windows для выполнения команд
4. Сообщите о проблеме в репозиторий Cursor на GitHub

