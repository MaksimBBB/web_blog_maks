# Персональний блог

## Можливості
- Перегляд списку статей на головній сторінці.
- Перегляд окремої статті.
- Адмін-панель (`dashboard`) зі списком статей.
- Створення, редагування та видалення статей.
- Валідація полів статті (`title`, `content`, `date`).
- Аутентифікація через JWT + cookie (`auth_token`).
- Захист адмін-маршрутів middleware-ом.
- Контейнеризація через Docker і `docker compose`.

## Структура проєкту

```text
web_blog_go/
├─ articles/
│  └─ 
├─ handlers/
│  ├─ article_handlers.go
│  ├─ auth.go
│  └─ middleware/
│     └─ jwt_middleware.go
├─ model/
│  └─ article.go
├─ templates/
│  ├─ home.html
│  ├─ articlepage.html
│  ├─ dashboard.html
│  ├─ newArticle.html
│  └─ updateArticle.html
├─ main.go
├─ user.go
├─ user_test.go
└─ .env
```

## Налаштування .env

Створіть файл `.env` у корені проєкту на основі `.env.example`:

```env
ADMIN_USERNAME=(your admin username)
ADMIN_PASSWORD=(your password)
JWT_SECRET=(your secret key)
```
## Запуск

1. Встановіть Go
2. У корені проєкту виконайте:

```bash
go mod tidy
go run main.go
```
Сервер буде доступний на:

`http://localhost:8080`

## Запуск через Docker

Проєкт містить:

- `Dockerfile` з multi-stage build.
- `docker-compose.yml` для запуску сервісу блогу.
- Named volume для збереження статей між перезапусками.

### Команди запуску

```
docker compose up --build
```

Після запуску застосунок буде доступний за адресою:

`http://localhost:8080`

## Маршрути

### Гостьова частина

- `GET /` - список статей
- `GET /articles` - список статей
- `GET /article/{id}` - сторінка статті

### Аутентифікація

- `GET /login` - форма логіну
- `POST /login` - вхід в адмінку, встановлення `auth_token`
- `GET /logout` - вихід (редірект на головну)

### Адмін (захищено)

- `GET /dashboard` - список статей з діями
- `GET /articles/new` - форма створення
- `POST /articles/new` - створення статті
- `GET /articles/update/{id}` - форма редагування
- `POST /articles/update/{id}` - оновлення статті
- `POST /articles/delete/{id}` - видалення статті

