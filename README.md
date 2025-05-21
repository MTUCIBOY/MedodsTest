# MedodsTest

## Описание
В данном сервисе авторизации было реализовано:
- Выдача пары токенов (access и refresh).
- Обновление пары токенов (access и refresh).
- Получение GUID пользователя.
- Деавторизация.
- Регестрация пользователя.

Проект использует:
- **Golang** (1.24.2)
- **chi** для маршрутизации.
- **PostgreSQL** для хранения данных.
- **JWT** для токенов.


## Структура маршрутов:
```go
POST /authTokens        // Получение пары токенов
POST /registrate        // Регистрация пользователя
GET  /guid              // Получение GUID (только аутентифицированные)
POST /deauthTokens      // Деавторизация
POST /updateTokens      // Обновление токенов (с проверкой IP/User-Agent)
```

## Установка и запуск

### 1. Запуск сервиса:
```bash
go run cmd/main.go
```
### 2. Запуск контейнера:
```bash
docker compose up -d
```

### 3. Документация API:
Откройте [http://localhost:8888/swagger/](http://localhost:8888/swagger/) для просмотра Swagger-документации.


##  Middleware
Для проверок данных были реализованы middleware:
- checkRefreshToken
- chechUserAgent
- checkUserIP
- expireTokenValidator
- tokenValidator

##  База данных
```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
	uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);
  
CREATE TABLE IF NOT EXISTS refresh_hashes (
	uuid UUID PRIMARY KEY,
	user_uuid UUID REFERENCES users(uuid),
	token_hash TEXT NOT NULL,
	is_active BOOL DEFAULT TRUE NOT NULL,
	created_at TIMESTAMP DEFAULT NOW() NOT NULL
);
```

##  Конфигурация
Настройки хранятся в `config.Config`:
```yaml
env: "local"
ttl_token: 10m
webhookURL: http://localhost:8888

http_server:
  address: "0.0.0.0:8888"
  read_timeout: 5s
  write_timeout: 5s
  idle_timeout: 10s
```
## Зависимости
- github.com/go-chi/chi/v5  v5.2.1
- github.com/golang-jwt/jwt/v5  v5.2.2
- github.com/google/uuid  v1.6.0
- github.com/ilyakaznacheev/cleanenv  v1.5.0
- github.com/jackc/pgx/v5  v5.7.5
- golang.org/x/crypto  v0.37.0