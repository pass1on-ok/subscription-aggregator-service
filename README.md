# Subscription Aggregator Service  

REST-сервис для агрегации данных об онлайн-подписках пользователей.  

## 📌 Функционал
- CRUDL-операции над записями о подписках:
  - **Название сервиса** (string)  
  - **Стоимость подписки** (int, рубли)  
  - **ID пользователя** (UUID)  
  - **Дата начала подписки** (месяц и год)  
  - **Дата окончания подписки** (опционально)  
- Подсчёт суммарной стоимости подписок за период  
  - фильтрация по `user_id`  
  - фильтрация по `service_name`  
- Миграции для PostgreSQL  
- Логирование операций  
- Конфигурация через `.env`  
- Swagger-документация  
- Запуск через **Docker Compose**

---

## ⚙️ Стек технологий
- **Go** 1.24  
- **PostgreSQL** 15  
- **Docker / Docker Compose**  
- **Swagger** (OpenAPI)  
- **Migrations** (через init-скрипты в `migrations/`)  

---

## 🚀 Запуск проекта

### 1. Клонировать репозиторий
```
git clone  https://github.com/pass1on-ok/subscription-aggregator-service.git
cd subscription-aggregator-service
```

### 2. Создать .env файл
```
APP_PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
```

### 3. Поднять сервис
```
docker-compose up --build
```

### 📂 Структура проекта
```
.
├── cmd/
│   └── main.go          # Точка входа
├── internal/
│   ├── handlers/        # HTTP-ручки
│   ├── models/          # Бизнес-логика и структуры
│   ├── repository/      # Работа с БД
│   └── service/         # Логика приложения
├── migrations/          # SQL-миграции
├── Dockerfile
├── docker-compose.yml
├── .env
└── README.md
```

### 📖 API (Swagger)

После запуска сервис доступен по адресу:
```
http://localhost:8080/swagger/index.html
```

### 📌 Примеры запросов
Создать подписку
```
curl -X POST http://localhost:8080//subscriptions \
-H "Content-Type: application/json" \
-d '{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}'
```

Получить список подписок
```
curl http://localhost:8080//subscriptions
```

Получить подписку по id
```
curl http://localhost:8080/subscriptions/{id}
```

Обновить подписку
```
curl -X PUT http://localhost:8080/subscriptions/{id} \
-H "Content-Type: application/json" \
-d '{
  "price": 500,
  "end_date": "12-2025"
}'
```

Удалить подписку
```
curl -X DELETE http://localhost:8080/subscriptions/{id}
```

Подсчитать сумму за период
```
curl "http://localhost:8080/subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&from=20
```
