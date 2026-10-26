# Taskflow API

REST API для управления задачами на Go.

## Технологии
- Go, Gin
- PostgreSQL, GORM
- Redis
- JWT
- Docker Compose
- Worker Pool

## Запуск

```bash
docker compose up -d --build
```

## API
- GET /health
- POST /auth/register
- POST /auth/login
- POST /tasks
- GET /tasks
- GET /tasks/:id
- PUT /tasks/:id
- DELETE /tasks/:id
