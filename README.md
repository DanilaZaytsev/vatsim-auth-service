### README.md

# VATSIM Auth Service

Микросервис авторизации пользователей через VATSIM SSO с сохранением данных в Yandex Database и выпуском собственного JWT токена.

## Возможности

* Авторизация через VATSIM SSO
* Сохранение пользователей в YDB
* Выпуск JWT токена (с ролью, email, страной и дивизионом)
* Middleware для проверки токена и ролей
* Ручка `/me` для получения информации о пользователе
* Ручка `/token` для получения токена из cookie
* Ручка `/auth/vatsim/login` и `/auth/vatsim/callback` для OAuth
* Ручка `/update-role` для обновления роли по CID (только для admin)
* Мониторинг `/admin/monitoring` (упрощённый)
* Health & Ready ручки

## API

| Method | Path                    | Description                       | Auth | Role  |
| ------ | ----------------------- | --------------------------------- | ---- | ----- |
| GET    | `/auth/vatsim/login`    | Перенаправление на VATSIM         | ❌    | -     |
| GET    | `/auth/vatsim/callback` | Обработка callback'а VATSIM       | ❌    | -     |
| GET    | `/me`                   | Текущий пользователь              | ✅    | любой |
| GET    | `/token`                | Получить JWT токен                | ✅    | любой |
| POST   | `/update-role`          | Обновить роль пользователя по CID | ✅    | admin |
| GET    | `/admin/monitoring`     | Примитивный мониторинг            | ✅    | admin |
| GET    | `/health`               | Проверка живости сервиса          | ❌    | -     |
| GET    | `/ready`                | Готовность к работе               | ❌    | -     |

## Запуск локально

```bash
cp .env.example .env
make run
```

## Переменные окружения

```env
YDB_DSN=grpcs://...
YDB_SA_KEY=./key.json
VATSIM_URL=https://auth-dev.vatsim.net
VATSIM_CLIENT_ID=...
VATSIM_CLIENT_SECRET=...
VATSIM_REDIRECT_URI=http://localhost:8080/auth/vatsim/callback
```

## Docker

```dockerfile

```

## GitHub Actions CI/CD (push в Yandex Container Registry)

```yaml
name: Docker Build and Push

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: docker/setup-buildx-action@v2
    - name: Login to Yandex Container Registry
      uses: docker/login-action@v2
      with:
        registry: cr.yandex
        username: oauth
        password: ${{ secrets.YC_OAUTH_TOKEN }}

    - name: Build and push
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: cr.yandex/your-registry-id/vatsim-auth-service:${{ github.sha }}
```

---

### swagger.yaml

```yaml
openapi: 3.0.0
info:
  title: VATSIM Auth Service
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /auth/vatsim/login:
    get:
      summary: Redirect to VATSIM login
      responses:
        '302':
          description: Redirect

  /auth/vatsim/callback:
    get:
      summary: Handle VATSIM callback
      responses:
        '200':
          description: JWT set in cookie

  /me:
    get:
      summary: Get current user
      security:
        - bearerAuth: []
      responses:
        '200':
          description: User profile

  /token:
    get:
      summary: Return JWT from cookie
      responses:
        '200':
          description: Token string

  /update-role:
    post:
      summary: Update role by CID
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              properties:
                cid:
                  type: integer
                role:
                  type: string
      responses:
        '200':
          description: Role updated

  /admin/monitoring:
    get:
      summary: Admin monitoring endpoint
      security:
        - bearerAuth: []
      responses:
        '200':
          description: Monitoring info

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
```

---