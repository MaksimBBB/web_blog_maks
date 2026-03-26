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
├─ .dockerignore
├─ .env.example
├─ Dockerfile
├─ docker-compose.yml
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

## Docker Команди

Запуск застосунку:

```bash
docker compose up --build
```
Зупинка контейнерів:

```bash
docker compose down
```

Перезапуск із перевідтворенням контейнера:

```bash
docker compose up --build --force-recreate
```

Видалення контейнерів разом із volume:

```bash
docker compose down -v
```

## Маршрути

### Гостьова частина

- `GET /` - головна сторінка зі списком статей.
- `GET /articles` - список статей.
- `GET /article/{id}` - перегляд окремої статті.
- `GET /login` - сторінка входу.
- `POST /login` - авторизація адміністратора.
- `GET /logout` - вихід із системи.
- `GET /dashboard` - адмін-панель.
- `GET /articles/new` - форма створення статті.
- `POST /articles/new` - створення статті.
- `GET /articles/update/{id}` - форма редагування.
- `POST /articles/update/{id}` - оновлення статті.
- `POST /articles/delete/{id}` - видалення статті.

## Troubleshooting

### `.env file not found`
Якщо застосунок запущений у Docker, це може бути нормально: змінні середовища передаються через `docker compose`, а не обов'язково через фізичний `.env` усередині контейнера.

### `port is already allocated`
Означає, що порт `8080` уже зайнятий іншим процесом або контейнером. Зупиніть конфліктний процес або змініть зовнішній порт у [docker-compose.yml]

### Статті не відображаються в контейнері
Якщо volume був створений раніше порожнім, потрібно перестворити його:
```bash
docker compose down -v
docker compose up --build
```

### Docker не запускається або не видно контейнер
Переконайтеся, що Docker Desktop запущений і має доступ до Linux середовища.
