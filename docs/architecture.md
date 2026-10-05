[← Конфигурация](configuration.md) · [Back to README](../README.md) · [Тестирование →](testing.md)

# Архитектура

Проект следует многослойной архитектуре (Clean Architecture): CLI-команды — тонкий слой, бизнес-логика вынесена в use-case'ы, доступ к данным — в репозитории.

## Структура каталогов

```
main.go                          # точка входа: .env, БД, регистрация команд
cmd/                             # CLI-команды (cobra), подпакеты по имени команды
  sync/ months/ year/ dynamics/ capital/ accounts/ migrate/ list/ menu/
internal/
  app/                           # контейнер зависимостей, инициализация БД
  config/                        # конфигурация из переменных окружения
  model/                         # доменные модели
  services/zenmoney/             # HTTP-клиент ZenMoney API (DTO)
  usecase/                       # бизнес-логика (Sync, Month, Year, Capital, Dynamics)
  adapter/
    db/                          # абстракция над GORM (DBServiceInterface)
    sqliterepo/zenrepo/
      accounts/                  # репозиторий счетов
      transactions/              # репозиторий транзакций
  dbinit/                        # AutoMigrate схемы БД
```

## Поток зависимостей

1. `main.go` создаёт `app.DB` (GORM + SQLite `zenmoney.db`), затем `app.Container` и `app.App`.
2. `Container` предоставляет фабрики `GetTransactionRepository()` и `GetAccountRepository()`.
3. CLI-команды (`cmd/...`) получают `*app.App`, извлекают репозитории/сервисы и создают use-case'ы «на лету».
4. Use-case'ы работают напрямую с репозиториями (`Month`, `Year`, `Accounts`, `Capital`), либо через `db.DBServiceInterface` — как `Sync`.
5. `Sync` вызывает `services/zenmoney` (HTTP к ZenMoney API), конвертирует DTO в модели и сохраняет через `DBServiceInterface`.

Внедрение зависимостей ручное, через конструкторы: `NewMonth(repo)`, `NewSync(db, api)`.

## Модели данных

| Модель | Таблица | Ключевые поля |
|--------|---------|---------------|
| `Transaction` | `transactions` | Id, Date, Income, Outcome, IncomeAccount, OutcomeAccount, TagIds, Deleted, Comment |
| `Account` | `accounts` | Id, Title, Balance, StartBalance, Instrument (FK на `instruments`) |
| `Instrument` | `instruments` | Id (PK), Title, ShortTitle, Symbol, Rate |
| `Tag` | `tags` | Id (PK), Title |
| `SyncState` | `sync_state` | ID, LastSyncedAt, ServerTimestamp, UpdatedAt |

Поле `TagIds` в `Transaction` хранит ID тегов строкой через запятую. Связанные счета (`InAccount`/`OutAccount`) и их валюты подгружаются через GORM `Preload`.

## Стиль кода

- Комментарии и сообщения — на русском; имена идентификаторов — на английском.
- Интерфейсы описываются в файле реализации и заканчиваются на `Interface` (`RepositoryInterface`, `DBServiceInterface`, `ApiInterface`).
- DTO имеют суффикс `Dto` (`MonthStatDto`, `AccountDto`).
- Ошибки оборачиваются через `fmt.Errorf("...: %w", err)`.

## Известные особенности

В корне `cmd/` остались устаревшие упрощённые версии команд (`cmd/sync.go`, `cmd/months.go` и т. д.). Актуальные определения находятся в подпакетах `cmd/<имя>/` и импортируются из `main.go`.

## Смотрите также

- [Тестирование](testing.md) — моки и структура тестов
- [Конфигурация](configuration.md) — инициализация БД и окружения
- [Участие в разработке](contributing.md) — соглашения проекта
