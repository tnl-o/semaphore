# Настройка валидации входных данных

## Установка зависимости

Для работы валидации необходимо добавить библиотеку `go-playground/validator`:

```bash
go get github.com/go-playground/validator/v10
go mod tidy
```

## Что было сделано

1. **Расширена функция `helpers.Bind`** - теперь она автоматически валидирует структуры с тегами `validate`
2. **Добавлены теги валидации** к основным структурам:
   - `db.User` - валидация username, name, email
   - `db.UserWithPwd` - валидация пароля (минимум 8 символов)
   - Структуры запросов в `api/login.go`, `api/auth.go`, `api/projects/users.go`

## Использование

Валидация работает автоматически при использовании `helpers.Bind()`. Просто добавьте теги `validate` к полям структуры:

```go
type MyRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
    Email    string `json:"email" validate:"required,email"`
    Age      int    `json:"age" validate:"required,min=18,max=120"`
}
```

## Доступные теги валидации

- `required` - поле обязательно
- `email` - валидный email адрес
- `min=n` - минимальная длина/значение
- `max=n` - максимальная длина/значение
- `alphanum` - только буквы и цифры
- `alpha` - только буквы
- `numeric` - только цифры
- `url` - валидный URL
- `oneof=value1 value2` - одно из указанных значений

## Защита от SQL-инъекций и XSS

Валидация входных данных помогает защитить от:
- SQL-инъекций - проверка формата и длины входных данных
- XSS атак - валидация формата данных перед обработкой
- Некорректных данных - раннее обнаружение проблем

