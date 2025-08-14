# 💰 Expense Tracker API

REST API для управления личными финансами, построенный на Go с использованием PostgreSQL.

## ✨ Возможности

- 👤 **Аутентификация пользователей** (регистрация/вход с JWT)
- 📝 **Управление транзакциями** (доходы/расходы)
- 🏷️ **Категории расходов** (создание пользовательских категорий)
- 🔍 **Фильтрация и поиск** (по типу, категории, датам)
- 📊 **Аналитика** (сводка доходов/расходов за период)
- 📄 **Пагинация** результатов
- 🔐 **Безопасность** (bcrypt для паролей, проверка прав доступа)

## 🛠️ Технологии

- **Backend:** Go 1.24
- **База данных:** PostgreSQL 15
- **Аутентификация:** JWT
- **Тестирование:** testify, mockery
- **Контейнеризация:** Docker, Docker Compose
- **CI/CD:** GitHub Actions

## 🚀 Быстрый старт

### Предварительные требования

- [Docker](https://docs.docker.com/get-docker/) и [Docker Compose](https://docs.docker.com/compose/install/)
- [Go 1.24+](https://golang.org/dl/) (для локальной разработки)

### Запуск с Docker Compose

```bash
# Клонируем репозиторий
git clone https://github.com/ViktorOHJ/expense-tracker.git
cd expense-tracker

# Запускаем приложение
docker-compose up --build

# Приложение будет доступно на http://localhost:8080
```

### Локальная разработка

```bash
# Запускаем PostgreSQL
docker run --name postgres-expense \
  -e POSTGRES_USER=expense_user \
  -e POSTGRES_PASSWORD=expense_pass \
  -e POSTGRES_DB=expense_tracker \
  -p 5432:5432 -d postgres:15

# Создаем .env файл
cat > .env << EOF
DB_URL=postgres://expense_user:expense_pass@localhost:5432/expense_tracker?sslmode=disable
JWT_SECRET=your-super-secret-jwt-key-here
PORT=8080
EOF

# Устанавливаем зависимости
go mod download

# Запускаем приложение
go run cmd/app/main.go
```

## 📚 API Документация

### Аутентификация

#### Регистрация
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Вход
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Транзакции

> 🔐 Все эндпоинты требуют заголовок `Authorization: Bearer <token>`

#### Создание транзакции
```http
POST /transactions
Content-Type: application/json
Authorization: Bearer <your-jwt-token>

{
  "is_income": false,
  "amount": 25.50,
  "category_id": 1,
  "note": "Обед в кафе"
}
```

#### Получение транзакций
```http
GET /transactions?limit=10&offset=1&type=false&category_id=1&from=2024-01-01&to=2024-12-31
Authorization: Bearer <your-jwt-token>
```

**Параметры запроса:**
- `limit` (обязательный) - количество записей
- `offset` (обязательный) - номер страницы (начиная с 1)
- `type` - тип транзакции (`true` для доходов, `false` для расходов)
- `category_id` - ID категории
- `from` - начальная дата (YYYY-MM-DD)
- `to` - конечная дата (YYYY-MM-DD)

#### Получение транзакции по ID
```http
GET /transaction/?id=1
Authorization: Bearer <your-jwt-token>
```

#### Удаление транзакции
```http
DELETE /transaction/?id=1
Authorization: Bearer <your-jwt-token>
```

### Категории

#### Создание категории
```http
POST /categories
Content-Type: application/json
Authorization: Bearer <your-jwt-token>

{
  "name": "Еда",
  "description": "Расходы на питание"
}
```

### Аналитика

#### Сводка за период
```http
GET /summary?from=2024-01-01&to=2024-12-31
Authorization: Bearer <your-jwt-token>
```

**Ответ:**
```json
{
  "message": "Summary retrieved successfully",
  "data": {
    "total_income": 5000.00,
    "total_expense": 3000.00,
    "balance": 2000.00
  }
}
```

## 🧪 Тестирование

```bash
# Запуск тестов
go test -v ./...

## 🏗️ Архитектура проекта

expense-tracker/
├── cmd/app/                 # Точка входа приложения
├── pkg/
│   ├── api/                 # HTTP handlers и middleware
│   │   └── handler_test/    # Тесты для handlers
│   ├── auth/                # JWT и работа с паролями
│   ├── db/                  # Слой работы с БД
│   ├── mocks/               # Моки для тестирования
│   ├── models.go            # Структуры данных
│   └── models_auth.go       # Структуры для аутентификации
├── .github/workflows/       # CI/CD конфигурация
├── Dockerfile
├── compose.yaml
├── go.mod
└── go.sum

## 🔧 Конфигурация

Приложение использует переменные окружения:

| Переменная | Описание | Значение по умолчанию |
|------------|----------|----------------------|
| `DB_URL` | URL подключения к PostgreSQL | - |
| `JWT_SECRET` | Секретный ключ для JWT | - |
| `PORT` | Порт для запуска сервера | `8080` |

## 🚀 CI/CD

Проект использует GitHub Actions для:
- ✅ Автоматического тестирования на каждый push/PR
- 🏗️ Сборки Docker образа
- 🔍 Проверки качества кода

## 🤝 Разработка

### Структура базы данных

```sql
-- Пользователи
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Категории
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, user_id)
);

-- Транзакции
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    is_income BOOLEAN NOT NULL,
    amount NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    category_id INTEGER REFERENCES categories(id),
    user_id INTEGER REFERENCES users(id),
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Добавление новых фичей

1. Создайте feature branch: `git checkout -b feature/amazing-feature`
2. Добавьте тесты для новой функциональности
3. Реализуйте функциональность
4. Убедитесь, что все тесты проходят: `go test ./...`
5. Создайте Pull Request

## 📄 Лицензия

Этот проект создан в учебных целях.

## 🙋‍♂️ Автор

**Viktor** - [GitHub](https://github.com/ViktorOHJ)

---
