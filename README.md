# Chirpy

A simple Twitter clone built with Go, featuring user management, post creation, and real-time interactions.

## About

Chirpy is a lightweight social media platform inspired by Twitter. It allows users to create accounts, post updates ("chirps"), and interact with other users' content. The project is built with Go and uses a PostgreSQL database for data persistence.

## Features

- **User Management**: Create accounts and manage user profiles
- **Chirps**: Post and view chirps (tweets) from users
- **Authentication**: Secure user authentication with JWT tokens
- **Database Migrations**: Managed with Goose for schema versioning
- **RESTful API**: Clean API endpoints for all operations

## Prerequisites

Before you can run Chirpy, ensure you have the following installed on your system:

- **Go** (version 1.21 or higher) - [Install Go](https://golang.org/doc/install)
- **PostgreSQL** (version 12 or higher) - [Install PostgreSQL](https://www.postgresql.org/download/)
- **Goose** (for database migrations) - Install with: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/adavidschmidt/Chirpy.git
cd Chirpy
```

### 2. Install Go Dependencies

```bash
go mod download
```

### 3. Set Up the Database

Create a PostgreSQL database for Chirpy:

```bash
createdb chirpy
```

### 4. Run Database Migrations

Use Goose to apply migrations:

```bash
goose -dir sql/schema postgres "user=postgres password=yourpassword dbname=chirpy sslmode=disable" up
```

### 5. Build and Run

```bash
go build -o Chirpy
./Chirpy
```

The server will start and listen for incoming requests. Check the console for the default port (typically `8080`).

## API Endpoints

For a complete description of all API endpoints, request methods, parameters, and response formats, see [API.md](./API.md).

## Project Structure

```
Chirpy/
├── sql/                 # Database migrations and schemas
├── internal/            # Internal packages
├── handler_chirps.go   # Chirps API handlers
├── handler_users.go    # Users API handlers
├── handler_tokens.go   # Token/authentication handlers
├── handler_webhooks.go # Webhook handlers
├── main.go             # Application entry point
├── go.mod              # Go module definition
└── README.md           # This file
```

## Configuration

Configure your database connection by setting environment variables or updating the connection string in the code:

```
DATABASE_URL=postgres://user:password@localhost:5432/chirpy
```

## Development

To run tests:

```bash
go test ./...
```

## License

This project is provided as-is for educational purposes.

## Contributing

Feel free to fork this repository and submit pull requests for any improvements!