# Project Name

[Project Structure](STRUCTURE.md)

**Project Description:** This project is a simple slot game application.

## Task Overview and Status

The table below summarizes all tasks required for the slot game REST API, indicating their status.

| Section              | Task Description                                                                                          | Status     |
|----------------------|-----------------------------------------------------------------------------------------------------------|------------|
| User Management      | Register a new user using email and password (`POST /api/register`)                                      | Completed  |
| User Management      | Login with email and password, providing token-based authorization (`POST /api/login`)                   | Completed  |
| User Management      | Retrieve user profile and credit balance (`GET /api/profile`)                                            | Completed  |
| Wallet Management    | Deposit credits to the user's balance (`POST /api/wallet/deposit`)                                      | Completed  |
| Wallet Management    | Withdraw credits from the user's balance (`POST /api/wallet/withdraw`)                                  | Completed  |
| Game Logic           | Spin slot machine (`POST /api/slot/spin`), bet, and calculate result                                     | Completed  |
| Game Logic           | Implement payout calculation based on winning combinations (e.g., two or three identical symbols)       | Completed  |
| Game History         | Retrieve game history (`GET /api/slot/history`)                                                          | Completed  |
| Technical Requirements | RESTful API implemented using Go                                                                         | Completed  |
| Technical Requirements | Use JWT for securing endpoints                                                                          | Completed  |
| Technical Requirements | Persist user data, transactions, and game history using PostgreSQL                                     | Completed  |
| Technical Requirements | Write unit tests for key components of the system                                                      | Completed  |
| Technical Requirements | Document the API using Swagger                                                                         | Completed  |
| Bonus Features       | Dockerize the application for easy setup                                                                 | Completed  |
| Bonus Features       | Add rate-limiting to prevent abuse (e.g., excessive spins)                                              | Completed  |
| Bonus Features       | Implement retry mechanism for spin function in case of failure                                          | Completed  |

## 1. Environment Setup

To successfully deploy the project, you need to install and configure several environment components. Follow the steps below:

### 1.1 Prerequisites
- **Docker**: Docker version 27.3.1 is required.
- **Docker Compose**: Docker Compose version 1.29.2 is required.
- **Go**: Go version 1.22.5 is required.
- **Git**: for working with the repository. Make sure Git is installed: `git --version`

### 1.2 Cloning the Repository
Clone the project from GitHub:
```bash
$ git clone github.com/vadymlab/slot-game
$ cd slot-game
```

### 1.3 Setting Up Environment Variables
Rename the `.env.example` file to `.env` and update the following parameters:
```bash
$ cp .env.example .env
```
Fill in the parameters in `.env`:
- **POSTGRES_USER**=
- **POSTGRES_PASSWORD**=
- **POSTGRES_DB**=
- **POSTGRES_HOST**=
- **REDIS_URL**=

This `.env` file will be used both for running the database in Docker Compose and for running the built application in Docker Compose.

## 2. Running the Database with Docker Compose

### 2.1 Command to Start the Database and Redis
Run the following command to start the database and Redis:
```bash
$ docker compose up -d database redis
```

The `docker-compose.yml` file contains the definitions for the database (e.g., PostgreSQL) and Redis. Make sure the ports and other settings match your environment.

After starting the containers, you can check the status with:
```bash
$ docker compose ps
```

## 3. Applying Database Migrations

### 3.1 Running Migrations
To populate the database with the necessary schema, you need to apply migrations. There are two ways to do this:

1. **Run migrations using Docker Compose**:
   ```bash
   $ docker compose up --build migrations
   ```

2. **Run migrations locally** (after installing `migrate` locally). You can find instructions here: [golang-migrate/migrate](https://github.com/golang-migrate/migrate/tree/master).

## 4. Running the Application

The application can be run in two main ways, depending on your environment.

### 4.0 Game Rules and Limits
- **Spin Rate Limiting**: By default, users are allowed to perform one spin per second.
- **Spin Retry Logic**: If a user does not have sufficient funds during a spin (e.g., a deposit transaction has not yet been processed), the spin will be retried for up to 2 seconds, with retries occurring every 500 milliseconds.

### 4.1 Running Locally
If you want to run the application locally (e.g., for development):
```bash
$ go run main.go
```

You can pass configuration parameters when running the application. To see a list of available options, run:
```bash
$ go run main.go --help
```

Here is a list of configuration parameters:

| Parameter                            | Description                                                                                                                              |
| ------------------------------------ |------------------------------------------------------------------------------------------------------------------------------------------|
| `--log-level value`                  | Application log level. Options: PANIC, FATAL, ERROR, WARNING, INFO, DEBUG, TRACE (default: "DEBUG") [\$LOG_LEVEL]                        |
| `--log-json`                         | Enable JSON format for logs. Defaults to false. (default: false) [\$LOG_JSON]                                                            |
| `--postgres-user value`              | PostgreSQL database username (default: "test") [\$POSTGRES_USER, \$PG_USER]                                                              |
| `--postgres-password value`          | PostgreSQL database user password (default: "test") [\$POSTGRES_PASSWORD, \$PG_PASSWORD]                                                 |
| `--postgres-db value`                | PostgreSQL database name (default: "node_art_slot_games") [\$POSTGRES_DB, \$PG_DB]                                                       |
| `--postgres-schema value`            | PostgreSQL database schema (default: "public") [\$POSTGRES_SCHEMA, \$PG_SCHEMA]                                                          |
| `--postgres-host value`              | PostgreSQL database host address (default: "localhost:5432") [\$POSTGRES_HOST, \$PG_HOST]                                                |
| `--postgres-max-life-time value`     | Maximum lifetime of a PostgreSQL connection in milliseconds (default: "20") [\$POSTGRES_CONNECTION_MAX_LIFE_TIME]                        |
| `--postgres-max-connection value`    | Maximum number of open PostgreSQL connections (default: "300") [\$POSTGRES_MAX_OPEN_CONNECTION]                                          |
| `--postgres-log-mode`                | Enable or disable query logging in PostgreSQL (default: true) [\$POSTGRES_QUERY_LOGGING]                                                 |
| `--server-host value`                | API server host address (default: "0.0.0.0") [\$API_HOST]                                                                                |
| `--server-port value`                | API server port (default: 8000) [\$API_PORT]                                                                                             |
| `--server-max-header-size value`     | Maximum size of request headers in bytes (default: 262144) [\$API_MAX_HEADER_SIZE]                                                       |
| `--server-request-timeout value`     | Maximum duration for reading the entire request in seconds (default: 5) [\$API_REQUEST_TIMEOUT]                                          |
| `--server-response-timeout value`    | Maximum duration before timing out writes of the response in seconds (default: 5) [\$API_RESPONSE_TIMEOUT]                               |
| `--server-log-request`               | Enable or disable request logging (default: true) [\$LOG_REQUEST]                                                                        |
| `--server-jwt-secret value`          | JWT secret used for signing authentication tokens (default: "qi87x8Sd9KpQUuiOMP7gFMid3gRTQFjr") [\$JWT_SECRET]                           |
| `--server-jwt-secret-lifetime value` | JWT token lifetime in minutes (default: 60) [\$JWT_SECRET_LIFE_TIME]                                                                     |
| `--multiplier-three value`           | Multiplier for three matching symbols (default: 10) [\$MULTIPLIER_THREE]                                                                 |
| `--multiplier-two value`             | Multiplier for two matching symbols (default: 2) [\$MULTIPLIER_TWO]                                                                      |
| `--two-match-probability value`      | Probability for winning with two matching symbols (default: 0.3) [\$TWO_MATCH_PROBABILITY]                                               |
| `--three-match-probability value`    | Probability for winning with three matching symbols (default: 0.05) [\$THREE_MATCH_PROBABILITY]                                          |
| ` --rate-limit value`                | Rate limit for requests per second( 5 reqs/second: "5-S", 10 reqs/minute: "10-M", 100 reqs/hour: "100-H") (default: "1-S") [\$RATE_LIMIT] |
| `--redis-url value`                  | Redis connection URL (default: "redis://localhost:6379/0") [\$REDIS_URL]                                                                |
| `--help, -h`                         | Show help                                                                                                                                |

### 4.2 Running with Docker Compose
To run the application with Docker Compose:
```bash
$ docker compose up -d --build
```

The application is available by default on port 8000: [http://localhost:8000/](http://localhost:8000/).

## 5. Swagger Documentation

The application has Swagger documentation available at: [http://localhost:8000/swagger/index.html#/](http://localhost:8000/swagger/index.html#/).

## 6. Postman Collection

You can use the following Postman collection for testing the API endpoints: [Slot Game Postman Collection](https://orange-meadow-363583.postman.co/workspace/SlotGames~d32ea593-fda8-40c2-a8fd-210ce73e7b6a/collection/4620563-bc182869-439e-463d-993a-3772a9737bbe?action=share&creator=4620563).

## 7. Game Rules and Limits

- **Spin Rate Limiting**: By default, users are allowed to perform one spin per second.
- **Spin Retry Logic**: If a user does not have sufficient funds during a spin (e.g., a deposit transaction has not yet been processed), the spin will be retried for up to 2 seconds, with retries occurring every 500 milliseconds.

