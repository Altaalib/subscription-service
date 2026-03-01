# Subscription Service

REST API сервис для агрегации данных об онлайн подписках пользователей.

## Стек

- **Go** + **Echo** — HTTP сервер
- **PostgreSQL** — база данных
- **sqlx + pgx** — работа с БД
- **golang-migrate** — миграции
- **zap** — логирование
- **swagger** — документация
- **Docker Compose** — запуск

## Запуск

### Требования
- Docker
- Docker Compose

### Как запустить
```bash
git clone https://github.com/твой-юзернейм/subscription-service
cd subscription-service
docker-compose up --build
```

Сервис запустится на `http://localhost:8080`

## API

### Swagger UI
```
http://localhost:8080/swagger/index.html
```

### Эндпоинты

| Метод | URL | Описание |
|-------|-----|----------|
| POST | /api/v1/subscriptions | Создать подписку |
| GET | /api/v1/subscriptions | Получить все подписки |
| GET | /api/v1/subscriptions/:id | Получить подписку по ID |
| PUT | /api/v1/subscriptions/:id | Обновить подписку |
| DELETE | /api/v1/subscriptions/:id | Удалить подписку |
| GET | /api/v1/subscriptions/total | Получить суммарную стоимость |

### Пример создания подписки
```json
{
    "service_name": "Yandex Plus",
    "price": 400,
    "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
    "start_date": "07-2025"
}
```

### Пример запроса суммы
```
GET /api/v1/subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&service_name=Yandex%20Plus
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| SERVER_PORT | Порт сервера | 8080 |
| DB_HOST | Хост БД | postgres |
| DB_PORT | Порт БД | 5432 |
| DB_USER | Пользователь БД | postgres |
| DB_PASSWORD | Пароль БД | postgres |
| DB_NAME | Имя БД | subscriptions |
| DB_SSLMODE | SSL режим | disable |