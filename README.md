# FinanceAppAPI

REST API для банковского сервиса с поддержкой регистрации, аутентификации, управления счетами, картами, переводами, кредитами и аналитикой.

## Возможности

- Регистрация и аутентификация пользователей (JWT)
- Создание и управление банковскими счетами
- Генерация и просмотр виртуальных карт (шифрование PGP, HMAC, bcrypt)
- Переводы между счетами
- Кредитование: оформление, график платежей, автоматическое списание, штрафы
- Финансовая аналитика (доходы/расходы, кредитная нагрузка, прогноз баланса)
- Интеграция с ЦБ РФ (ключевая ставка через SOAP)
- Email-уведомления через SMTP
- Защита данных: шифрование, хеширование, HMAC

## Требования

- Go 1.23+
- PostgreSQL 17+
- [Gorilla Mux](https://github.com/gorilla/mux)
- [lib/pq](https://github.com/lib/pq)
- [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
- [logrus](https://github.com/sirupsen/logrus)
- [beevik/etree](https://github.com/beevik/etree)
- [go-mail/mail/v2](https://github.com/go-mail/mail)

## Установка

1. Клонируйте репозиторий:
    ```sh
    git clone https://github.com/yourusername/financeAppAPI.git
    cd financeAppAPI
    ```

2. Установите зависимости:
    ```sh
    go mod download
    ```

3. Настройте переменные окружения (опционально, значения по умолчанию заданы в [internal/config/config.go](internal/config/config.go)):
    - `DATABASE_URL` — строка подключения к PostgreSQL
    - `JWT_SECRET` — секрет для подписи JWT
    - `SMTP_FROM`, `SMTP_PASS`, `SMTP_HOST`, `SMTP_PORT` — параметры SMTP
    - `PGP_PRIVATE_KEY`, `PGP_PASSPHRASE` — для шифрования карт

4. Создайте базу данных и примените миграции:
    ```sh
    psql -U postgres -d findbgo1 -f migrations/001_init.sql
    ```

5. Запустите приложение:
    ```sh
    go run cmd/main.go
    ```

## Использование

### Регистрация и вход

- `POST /register` — регистрация пользователя  
  Тело запроса:  
  ```json
  {
    "username": "user1",
    "email": "user1@example.com",
    "password": "StrongP@ssw0rd!"
  }
  ```

- `POST /login` — аутентификация  
  Тело запроса:  
  ```json
  {
    "email": "user1@example.com",
    "password": "StrongP@ssw0rd!"
  }
  ```
  Ответ:  
  ```json
  { "token": "..." }
  ```

### Защищённые эндпоинты (требуется JWT в заголовке Authorization: Bearer ...)

- `POST /api/accounts` — создать счет
- `POST /api/cards` — выпустить карту
- `GET /api/cards?account_id=1` — получить карты по счету
- `POST /api/transfers` — перевод между счетами
- `POST /api/accounts/deposit` — пополнение счета
- `POST /api/accounts/withdraw` — списание со счета
- `POST /api/credits` — оформить кредит
- `GET /api/credits/{creditId}/schedule` — график платежей по кредиту
- `GET /api/accounts/{accountId}/predict?days=30` — прогноз баланса
- `GET /api/analytics` — аналитика по счету

### Пример запроса

```sh
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","email":"user1@example.com","password":"StrongP@ssw0rd!"}'
```

## Примечания

- Для отправки email используйте тестовый SMTP-сервер (например, [MailHog](https://github.com/mailhog/MailHog)).
- Для работы с картами требуется PGP-ключ (см. переменные окружения).
- Прогноз баланса ограничен 365 днями.
- Все ответы возвращаются в формате JSON.

## Структура проекта

- `cmd/main.go` — точка входа
- `internal/app/` — инициализация приложения
- `internal/handlers/` — HTTP-обработчики
- `internal/services/` — бизнес-логика
- `internal/repositories/` — работа с БД
- `internal/models/` — структуры данных
- `internal/utils/` — утилиты (шифрование, валидация)
- `internal/config/` — конфигурация
- `migrations/` — SQL-миграции

---
