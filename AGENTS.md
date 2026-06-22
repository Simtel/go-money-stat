# AGENTS.md

## Обзор проекта

**Zenmoney CLI** (`money-stat`) — CLI-приложение на Go 1.22 для взаимодействия с
API [Zenmoney.ru](https://api.zenmoney.ru/v8/diff/). Позволяет синхронизировать финансовые данные (транзакции, счета,
валюты, теги) в локальную базу SQLite и анализировать их через терминал: просмотр транзакций за месяц/год, отчёты по
капиталу и остаткам по счетам.

## Технологический стек

- **Язык:** Go 1.22 (toolchain go1.22.2, модуль `money-stat`)
- **CLI-фреймворк:** [cobra](https://github.com/spf13/cobra) v1.8.1
- **ORM:** [GORM](https://gorm.io) v1.25.12 + драйвер SQLite (`gorm.io/driver/sqlite`)
- **Терминальный UI:** [pterm](https://github.com/pterm/pterm) v0.12.80 — таблицы, спиннеры, прогресс-бары
- **Конфигурация:** [godotenv](https://github.com/joho/godotenv) v1.5.1 — загрузка `.env`
- **Тестирование:** [testify](https://github.com/stretchr/testify) v1.8.4 + [gomock](https://github.com/golang/mock)
  v1.6.0 + go-sqlmock v1.5.2

## Сборка и запуск

```bash
# Установка зависимостей
go mod tidy

# Сборка бинарника
go build -o money-stat .

# Запуск (после копирования .env.example в .env и заполнения ZENMONEY_TOKEN)
go run . <команда>

# Основные команды
go run . sync                   # инкрементальная синхронизация с ZenMoney
go run . sync --full            # полная синхронизация (сброс и перезагрузка)
go run . months current         # транзакции за текущий месяц
go run . months previous        # транзакции за прошлый месяц
go run . year 2025              # отчёт доходов/расходов за год
go run . capital 2025           # помесячный капитал за год
go run . accounts               # список счетов с балансами
go run . migrate init           # инициализация/миграция БД
```

## Тестирование

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов конкретного пакета
go test ./internal/usecase/...

# С флагом verbose
go test -v ./...
```

CI настроен через GitHub Actions (`.github/workflows/tests.yml`): на push/PR в master запускаются `go test ./...` и
`golangci-lint` v1.59.

## Архитектура

Проект следует многослойной архитектуре (Clean Architecture):

```
main.go                          # точка входа: загрузка .env, инициализация БД, регистрация команд
cmd/                             # определения CLI-команд (cobra)
  sync/sync.go, months/months.go, year/year.go, accounts/accounts.go,
  capital/capital.go, migrate/migrate.go
  sync.go, months.go, ...        # дублирующие/устаревшие версии в корне cmd/
internal/
  app/                           # контейнер зависимостей (Container), инициализация БД (DB)
  config/                        # конфигурация из переменных окружения
  model/                         # доменные модели (Transaction, Account, Instrument, Tag, SyncState)
  services/zenmoney/             # HTTP-клиент ZenMoney API (DTO ответов, запрос Diff/DiffSince)
  usecase/                       # бизнес-логика (Sync, Month, Year, Accounts, Capital)
  adapter/
    db/                          # абстракция-обёртка над GORM (DBServiceInterface)
    sqliterepo/zenrepo/
      accounts/                  # репозиторий счетов (RepositoryInterface)
      transactions/              # репозиторий транзакций (RepositoryInterface)
  dbinit/                        # инициализация схемы БД (AutoMigrate)
```

### Поток зависимостей

1. `main.go` создаёт `app.DB` (GORM + SQLite `zenmoney.db`), затем `app.Container` и `app.App`.
2. `Container` предоставляет `GetTransactionRepository()` и `GetAccountRepository()` — фабрики репозиториев.
3. CLI-команды (`cmd/...`) получают `*app.App`, извлекают нужные репозитории/сервисы и создают use-case-объекты «на
   лету».
4. Use-case'ы (`usecase/`) работают либо напрямую с репозиториями (`Month`, `Year`, `Accounts`, `Capital`), либо через
   `db.DBServiceInterface` (абстракция над GORM) — как `Sync`.
5. Синхронизация (`Sync`) вызывает `services/zenmoney` (HTTP к ZenMoney API), конвертирует DTO в модели и сохраняет
   через `db.DBServiceInterface`.

### Модели данных

| Модель      | Таблица      | Ключевые поля                                                                      |
|-------------|--------------|------------------------------------------------------------------------------------|
| Transaction | transactions | Id, Date, Income, Outcome, IncomeAccount, OutcomeAccount, TagIds, Deleted, Comment |
| Account     | accounts     | Id, Title, Balance, StartBalance, Instrument (FK на instruments)                   |
| Instrument  | instruments  | Id (PK), Title, ShortTitle, Symbol, Rate                                           |
| Tag         | tags         | Id (PK), Title                                                                     |
| SyncState   | sync_state   | ID, LastSyncedAt, ServerTimestamp, UpdatedAt                                       |

Важно: поле `TagIds` в `model.Transaction` хранит ID тегов строкой через запятую, а связанные `Tag` и `Account` (
`InAccount`/`OutAccount`) подгружаются отдельно через GORM Preload/Joins только в репозитории счетов. В репозитории
транзакций связанные сущности **не** подгружаются.

### Двойные реализации команд

В проекте существуют два набора определений команд:

- `cmd/sync/sync.go`, `cmd/months/months.go`, ... — **актуальные**, импортируются из `main.go`
- `cmd/sync.go`, `cmd/months.go`, ... — устаревшие/упрощённые версии в пакете `cmd`

При добавлении новых команд используйте подпакет `cmd/<имя>/`. Не изменяйте файлы в корне `cmd/` без явной
необходимости.

## Стиль кода

- **Язык комментариев и сообщений:** русский (интерфейс, логи, описания команд)
- **Имена идентификаторов:** английский (переменные, функции, типы)
- **Именование пакетов:** строчные буквы, короткие (`usecase`, `zenmoney`, `sqliterepo`)
- **Интерфейсы:** описываются в том же файле, где определена структура-реализация; имя заканчивается на `Interface` (
  `RepositoryInterface`, `DBServiceInterface`, `ApiInterface`)
- **DTO:** суффикс `Dto` (например, `MonthStatDto`, `AccountDto`)
- **Внедрение зависимостей:** ручное, через конструктор (`NewMonth(repo)`, `NewSync(db, api)`)
- **Обработка ошибок:** стандартная для Go — возврат `error`, оборачивание через `fmt.Errorf("...: %w", err)`

## Конфигурация и безопасность

- Обязательная переменная окружения: `ZENMONEY_TOKEN` — токен доступа к ZenMoney API
- Задаётся в файле `.env` (не коммитится, есть в `.gitignore`)
- Шаблон: `.env.example` (содержит только ключ с placeholder-значением)
- База данных SQLite (`zenmoney.db`) также в `.gitignore`, создаётся автоматически при первом запуске
- **Никогда не коммитьте токены или реальные значения из `.env`**

## Принятые соглашения при разработке

1. Миграции БД выполняются через `dbinit.InitializeDB()` с использованием `AutoMigrate` — новые поля моделей добавляются
   автоматически.
2. Синхронизация по умолчанию инкрементальная — флаг `--full` для полной перезагрузки.
3. Для операций upsert используется паттерн: поиск существующей записи через `First`, затем `Save` или `Create`.
4. В CLI-выводе активно используется `pterm` для форматированных таблиц и цветового оформления (красный для
   отрицательных значений, зелёный для положительных).
5. Кэширование в `Capital` реализовано в памяти с TTL 5 минут.
