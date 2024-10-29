Here’s a structured project description for your README.md in English:

---

# Slot Game Backend

## Project Structure

The project is organized into three primary layers, enhancing modularity and code readability:

### 1. Controllers (`controller`)

Controllers manage HTTP requests, handle routing, input validation, and response formatting. They interact with the service layer to process requests, passing the required data to services and returning formatted responses to the client. Controllers are designed to be minimal, focusing on delegating business logic to the service layer.

### 2. Services (`service`)

The service layer encapsulates the core business logic of the application. It processes the main operations, manages data manipulation, and interacts with repositories for database access. Services ensure transaction consistency by using a Unit of Work pattern, coordinating data operations across multiple repositories when necessary.

### 3. Repositories (`repository`)

Repositories provide direct data access and abstract database operations, allowing the service layer to interact with data in a clean, decoupled way. Each repository is responsible for managing data retrieval, insertion, updates, and deletion for specific models, such as `User` or `Spin`.

## Transaction Management with Unit of Work

The project implements the Unit of Work pattern across all layers to ensure transactional consistency. By using [gorm-unit-of-work](https://github.com/public-forge/gorm-unit-of-work) for PostgreSQL, transactions are managed seamlessly across service and repository calls, allowing rollback and commit operations to be handled in a structured way. This ensures that related database operations are either all successfully applied or all reverted, maintaining data integrity.

## Dependency Injection with `go.uber.org/fx`

For dependency injection, the project uses [go.uber.org/fx](https://pkg.go.dev/go.uber.org/fx), which simplifies application startup, dependency graph construction, and lifecycle management. This framework enables controllers, services, and repositories to be injected at runtime, ensuring a loosely coupled architecture that’s easy to test and extend.

## Application Startup with `github.com/urfave/cli/v2`

Application startup and configuration management are handled by [github.com/urfave/cli/v2](https://github.com/urfave/cli/v2). This library allows for flexible configuration of application parameters through both command-line arguments and environment variables. CLI flags are used to control API server settings, such as host, port, JWT secrets, and other key parameters, providing a consistent and manageable approach to configuration.