# API Documentation for Chirpy

## Chirps Endpoints

### Get Chirps
- **Endpoint:** `GET /chirps`
- **Description:** Retrieve a list of chirps.
- **Response:**
  - 200 OK
  - 404 Not Found

### Create Chirp
- **Endpoint:** `POST /chirps`
- **Description:** Create a new chirp.
- **Request Body:
  ```json
  {
    "content": "string"
  }
  ```
- **Response:**
  - 201 Created
  - 400 Bad Request

## Users Endpoints

### Get Users
- **Endpoint:** `GET /users`
- **Description:** Retrieve a list of users.
- **Response:**
  - 200 OK
  - 404 Not Found

### Create User
- **Endpoint:** `POST /users`
- **Description:** Register a new user.
- **Request Body:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response:**
  - 201 Created
  - 400 Bad Request

## Authentication Endpoints

### Login
- **Endpoint:** `POST /auth/login`
- **Description:** Authenticate a user and receive a token.
- **Request Body:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response:**
  - 200 OK
  - 401 Unauthorized

### Logout
- **Endpoint:** `POST /auth/logout`
- **Description:** Logout the authenticated user.
- **Response:**
  - 200 OK
  - 401 Unauthorized

## Health Check Endpoints

### Health Check
- **Endpoint:** `GET /health`
- **Description:** Check the health of the API.
- **Response:**
  - 200 OK
  - 503 Service Unavailable

## Metrics Endpoints

### Get Metrics
- **Endpoint:** `GET /metrics`
- **Description:** Retrieve system metrics.
- **Response:**
  - 200 OK
  - 500 Internal Server Error