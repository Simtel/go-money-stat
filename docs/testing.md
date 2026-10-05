[← Архитектура](architecture.md) · [Back to README](../README.md) · [Участие в разработке →](contributing.md)

# Тестирование

Проект покрыт модульными тестами на базе [testify](https://github.com/stretchr/testify), с моками через [gomock](https://github.com/golang/mock) и [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock).

## Запуск тестов

```bash
# Все тесты
go test ./...

# Конкретный пакет
go test ./internal/usecase/...

# Подробный вывод
go test -v ./...
```

## Что покрыто

| Пакет | Тесты |
|-------|-------|
| `internal/model` | Методы `Account`, `Transaction` (форматирование суммы, типы транзакций) |
| `internal/usecase` | `Month`, `Year`, `Accounts`, `Capital`, `Sync` |
| `internal/services/zenmoney` | Клиент API, DTO запросов/ответов |
| `internal/adapter/sqliterepo/zenrepo` | Репозитории счетов и транзакций |
| `internal/app`, `internal/config` | Контейнер зависимостей и конфигурация |

## Моки

- `internal/adapter/sqliterepo/zenrepo/*/mocks/` — моки репозиториев, сгенерированные gomock.
- `internal/app/mock_db.go` — мок интерфейса БД.
- `go-sqlmock` — имитация SQL-запросов в тестах репозиториев.

## CI

GitHub Actions (`.github/workflows/tests.yml`) запускается на push и pull request в `master`. Две задачи:

1. **test** — `go mod tidy` и `go test ./...` на Go 1.22.
2. **golangci** — `golangci-lint` версии v1.59.

Перед тестами CI копирует `.env.example` в `.env`, чтобы переменные окружения были доступны.

## Смотрите также

- [Архитектура](architecture.md) — слои, которые покрывают тесты
- [Участие в разработке](contributing.md) — запуск линтера перед PR
