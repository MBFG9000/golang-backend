# Practice 4 API

Go REST API for managing users, containerized with Docker.

## Requirements

- Docker
- Docker Compose

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd <repository-name>
```

### 2. Create `.env` file

Create a `.env` file in the **root folder** of the project:

```
HOST="db"
DB_USERNAME="postgres"
PASSWORD="postgres"
DATABASE_NAME="mydb"
SSL_MODE="disable"
API_KEY_HEADER="X-API-KEY"
VALID_API_KEY="secretkey"
```

### 3. Create `postgres.env` file

Create a `postgres.env` file in the **root folder** of the project:

```
POSTGRES_PASSWORD=postgres
POSTGRES_USER=postgres
POSTGRES_DB=mydb
```

### 4. Run the application

```bash
docker compose up --build
```

The app will be available at `http://localhost:8080`.

> ⚠️ The application waits for the database to be healthy before starting.

## Local development seeding

You can generate development data with the Go seeder command:

```bash
go run ./cmd/seeder --users 200 --friends 8 --truncate --seed 42
```

Flags:

- `--users` number of users to generate (default: `100`)
- `--friends` max number of friend links per user (default: `5`)
- `--truncate` clear `users` and `user_friends` before seeding
- `--seed` random seed for deterministic generation

If you run inside Docker:

```bash
docker compose up -d --build
docker compose exec backend /seeder --users 200 --friends 8 --truncate --seed 42
```

## Authentication

All API endpoints are protected by API Key authentication.

Add the API key to every request via header:

```
X-API-KEY: secretkey
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
curl -H "X-API-KEY: secretkey" http://localhost:8080/users
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET    | `/users` | Get all users if no query params; returns paginated response if query params are provided (`page`, `page_size`, filters, sorting) |
| GET    | `/users/{id}` | Get user by UUID |
| GET    | `/users/common-friends?user_id_1=<uuid>&user_id_2=<uuid>` | Get common friends between two users |
| POST   | `/users` | Create user |
| PUT    | `/users/{id}` | Update user by UUID |
| DELETE | `/users/{id}` | Delete user by UUID |
