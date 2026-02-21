# Practice 3 API

Go REST API for managing users

## Requirements

- Go 1.22+
- PostgreSQL

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd <repository-name>
```

### 2. Create `.env` file

Create a `.env` file in the **root folder** of the project:

```
HOST="localhost"
DB_USERNAME="postgres"
PASSWORD="ADMIN"
DATABASE_NAME="go_db"
SSL_MODE="disable"
API_KEY_HEADER="X-API-KEY"
VALID_API_KEY="secret12345"
```

> ⚠️ The `.env` file must be in the root folder, otherwise the app won't find it.

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the application

Always run from the **root folder**:

```bash
go run cmd/api/main.go
```

> ⚠️ Do not run from inside `cmd/api/` — migrations and `.env` are resolved relative to the root.

## Authentication

All API endpoints are protected by API Key authentication.

Add the API key to every request via header:
```
X-API-KEY: secret12345
```

If the header is missing or invalid, the API will return:
```json
{
  "error": "unauthorized"
}
```

with status code `401 Unauthorized`.

### Example request
```bash
curl -H "X-API-KEY: secret12345" http://localhost:8080/users
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/users` | Get all users |
| GET    | `/users/{id}` | Get user by ID |
| POST   | `/users` | Create user |
| PUT    | `/users/{id}` | Update user |
| DELETE | `/users/{id}` | Delete user |