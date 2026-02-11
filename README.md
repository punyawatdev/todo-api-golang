# Todo API - Go Project Structure

## 📁 Project Structure (Standard Go Layout)

```
todo-api/
├── cmd/
│   └── api/
│       └── main.go           # Entry point
├── internal/
│   ├── database/
│   │   └── postgres.go       # Database connection
│   ├── todo/
│   │   ├── model.go          # Domain model
│   │   ├── repository.go     # Interface
│   │   ├── postgres_repo.go  # PostgreSQL implementation
│   │   ├── memory_repo.go    # In-memory implementation
│   │   ├── service.go        # Business logic
│   │   └── handler.go        # HTTP handlers / same Controller
│   └── server/
│       └── server.go         # HTTP server setup
├── migrations/
│   └── 001_create_todos.sql
├── .env
├── go.mod
└── go.sum
```

## 🚀 How to Run

### 1. Setup PostgreSQL

```bash
# Using Docker
docker run --name postgres-todo \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=tododb \
  -p 5432:5432 \
  -d postgres:16-alpine

# Or install PostgreSQL locally
```

### 2. Run Migrations

```bash
# Connect to PostgreSQL
psql -h localhost -U postgres -d tododb

# Run migration
\i migrations/001_create_todos.sql

# Or Docker Desktop
docker exec -i postgres_container psql -U postgres -d tododb < migrations/001_create_todos.sql
```

### 3. Install Dependencies

```bash
go mod init todo-api
go get github.com/jackc/pgx/v5
go get github.com/joho/godotenv
go mod tidy
```

### 4. Create .env file

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 5. Run the API

```bash
go run cmd/api/main.go
```

---

## 📝 Test API

```bash
# Create
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Docker Basics","description":"Understand containers, images, and docker-compose."}'

# List all
curl http://localhost:8080/todos

# Get one
curl http://localhost:8080/todos/1 

# Update
curl -X PUT http://localhost:8080/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Learn Docker Basics","description":"Done!","completed":true}'

# Delete
curl -X DELETE http://localhost:8080/todos/1
```

---

