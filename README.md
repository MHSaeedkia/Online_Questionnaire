
# Online Questionnaire Project

An online questionnaire platform built with Go and Fiber. This project includes an API for creating and managing questionnaires, with support for hot reloading using `air` and API documentation using Swagger.

## Table of Contents

- [Features](#features)
- [Requirements](#requirements)
- [Setup Instructions](#setup-instructions)
- [Running the Project](#running-the-project)
    - [Using Air](#using-air-for-hot-reloading)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Environment Variables](#environment-variables)
- [Contributing](#contributing)

---

## Features

- RESTful API built with Go and Fiber
- Middleware for authentication and logging
- Swagger integration for API documentation
- PostgreSQL database for data persistence
- Hot reloading with `air`

---

## Requirements

- [Go 1.20+](https://golang.org/dl/)
- [PostgreSQL 14+](https://www.postgresql.org/download/)
- [Air](https://github.com/cosmtrek/air) for live reloading
- [Docker (Optional)](https://www.docker.com/) for containerized development

---

## Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/online_questionnaire.git
cd online_questionnaire
```

### 2. Install Dependencies

Ensure Go modules are enabled:

```bash
go mod tidy
```

### 3. Configure Environment Variables

Create a `.env` file in the root directory and add the following variables:

```env
APP_NAME=OnlineQuestionnaire
DEBUG=true

DB_USER=admin
DB_PASSWORD=admin123
DB_NAME=questionnaire
DB_HOST=localhost
DB_PORT=5432
DB_SSLMODE=disable

SERVER_HOST=localhost
SERVER_PORT=8080
```

### 4. Set Up the Database

Ensure PostgreSQL is running. Then, create the database:

```sql
CREATE DATABASE questionnaire;
```

Run migrations if applicable.

---

## Running the Project

### Using Air for Hot Reloading

1. Install `air` globally if not already installed:

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. Run the project with `air`:

   ```bash
   air
   ```

   This will enable hot reloading. Any changes made to `.go` files in the watched directories will automatically trigger a rebuild and restart the application.

### Without Air

You can also run the project directly with Go:

```bash
go run cmd/app/main.go
```

---

## API Documentation

Swagger is used to generate API documentation.

1. Ensure the application is running.
2. Open your browser and navigate to:

   ```
   http://localhost:8080/swagger/index.html
   ```

### Setting Up Swagger

1. Install the Swagger generator:

   ```bash
   go get -u github.com/swaggo/swag/cmd/swag
   ```

2. Generate Swagger documentation:

   ```bash
   swag init --generalInfo cmd/app/main.go
   ```

3. The generated files will appear in the `docs` directory.

---

## Project Structure

```
online_questionnaire/
├── cmd/
│   └── app/
│       └── main.go           # Application entry point
├── configs/                  # Configuration files
├── docs/                     # Swagger documentation
├── internal/
│   ├── db/                   # Database setup and migrations
│   ├── handlers/             # API handlers
│   ├── middlewares/          # Middleware functions
│   ├── models/               # Data models
│   ├── repositories/         # Database interaction layer
│   ├── services/             # Business logic
├── scripts/                  # Utility scripts for automation
├── tmp/                      # Temporary files (ignored in version control)
├── .air.toml                 # Air configuration
├── .env                      # Environment variables
├── go.mod                    # Go modules
└── README.md                 # Project documentation

```

---

## Environment Variables

The project requires the following environment variables:

| Variable        | Description                      | Example        |
|-----------------|----------------------------------|----------------|
| `APP_NAME`      | Name of the application          | OnlineQuestionnaire |
| `DEBUG`         | Enable debug mode (true/false)   | true           |
| `DB_USER`       | Database username                | admin          |
| `DB_PASSWORD`   | Database password                | admin123       |
| `DB_NAME`       | Database name                   | questionnaire  |
| `DB_HOST`       | Database host                   | localhost      |
| `DB_PORT`       | Database port                   | 5432           |
| `DB_SSLMODE`    | Database SSL mode               | disable        |
| `SERVER_HOST`   | Server host                     | localhost      |
| `SERVER_PORT`   | Server port                     | 8080           |

---

## Contributing

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Push to the branch.
5. Open a pull request.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

Let me know if you need additional adjustments or explanations!