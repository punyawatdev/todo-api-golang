# 📝 Todo API (EXAMPLE)

> A robust, production-ready Todo API built with Go, following Clean Architecture principles. This project is fully containerized and features automated database migrations and graceful shutdown handling.

## 🛠 Tech Stack
- Language: Go 1.25+
- Database: PostgreSQL 15+
- Containerization: Docker & Docker Compose
- Architecture: Handler -> Service -> Repository (Interface-based)

## 📁 Project Structure (Standard Go Layout)

```
todo-api/
├── cmd/api/                  # Application Entry Point (main.go)
├── internal/
│   ├── database/             # DB Connection & Auto-migration Logic
│   │   └── postgres.go       # Database file
│   ├── todo/
│   │   ├── model.go          # Domain model file
│   │   ├── repository.go     # Interface repository file
│   │   ├── postgres_repo.go  # PostgreSQL implementation file
│   │   ├── memory_repo.go    # In-memory implementation file
│   │   ├── service.go        # Business logic file
│   │   └── handler.go        # HTTP handlers / same Controller file
│   └── server/               # HTTP Server & Graceful Shutdown Setup
│       └── server.go         # HTTP Server file
├── migrations/               # SQL Schema Definitions
│   └── 001_create_todos.sql. # Mock data file
├── docker-compose.yml        # Container Orchestration
├── Dockerfile                # Multi-stage Docker Build
├── .env                      # Environment Variables file
├── go.mod
└── go.sum
```

## 🚀 Quick Start (Recommended)

> The easiest way to get the API running is using Docker Compose. This will automatically set up the PostgreSQL database and the Go application.

### 1. Prepare Environment Variables

```bash
cp .env.example .env  # Configure your DB_USER, DB_PASSWORD, etc.
```

### 2. Spin up the Containers

```bash
docker-compose up --build
```
The system will automatically connect to the database and run SQL migrations.

### 3. Access the API

- API Endpoint: ```http://localhost:8080```
- Database Host: ```localhost:5432```


## 💻 Local Development

> To run the application locally without Docker (ensure you have a PostgreSQL instance running):

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Run Application

```bash
go run cmd/api/main.go
```

## 🧪 API Endpoints
| Method | Endpoint | Description |
| :--- | :--- | :--- |
| ```GET``` | ```/todos``` | Fetch all todo items | 
| ```POST``` | ```/todos``` | Create a new todo | 
| ```GET``` | ```/todos/:id``` | Get details of a specific todo | 
| ```PUT``` | ```/todos/:id``` | Update an existing todo | 
| ```DELETE``` | ```/todos/:id``` | Delete a todo | 

### Example: Create a Todo

```bash
curl -X POST http://localhost:8080/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"Master Go Docker","description":"Implement multi-stage builds"}'
```

## 🛡 Key Features
- ✅ Clean Architecture: Strict separation of concerns making the code testable and maintainable.

- ✅ Graceful Shutdown: Ensures all active requests are completed before the server stops, preventing data corruption.

- ✅ Automated Migrations: Database schemas are checked and updated automatically on startup.

- ✅ Dependency Injection: Easily switch between PostgreSQL and In-memory repositories for testing.

- ✅ Multi-stage Docker Build: Optimized images using Alpine Linux for a smaller footprint and better security.

## 🔧 Environment Configuration

| Variable | Description | Default |
| :--- | :--- | :--- |
| ```DB_HOST``` | Database host address | ```db``` (for Docker) |
| ```DB_PORT``` | Database port | ```5432``` |
| ```DB_USER``` | Database username | ```user``` |
| ```DB_PASSWORD``` | Database password | ```pass``` |
| ```DB_NAME``` | Database name | ```todo_db``` |
| ```SERVER_PORT``` | API server port | ```8080``` |