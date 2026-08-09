# User service

User management service for the Auth Platform.

The service is responsible for user lifecycle management, user data persistence, password management, and providing user information to other services through gRPC.

## Overview

**auth-user** is a backend microservice responsible for managing application users.

The service provides an isolated user management domain and exposes its functionality to other services through gRPC.

The service does not issue JWT tokens and does not make authorization decisions. Authentication is handled by **auth-auth**, while roles and permissions are managed by **auth-access**.

## Responsibilities

- User creation
- User retrieval
- User update
- User deletion
- User lookup by ID
- User lookup by credentials
- Password storage
- Password hashing
- User data persistence
- Providing user information to other services

## Repository Structure

```text
auth-user/
├── cmd/
│   └── main.go                     # Application entry point
├── configs/
│   ├── config-local.yml            # Local development configuration
│   └── config-prod.yml             # Production configuration
├── internal/
│   ├── config/
│   │   └── config.go               # Configuration parsing and loading
│   ├── domain/
│   │   ├── user.go                 # Domain models (User, etc.)
│   │   └── response.go             # Response DTOs
│   ├── handler/
│   │   ├── http-handler/
│   │   │   └── http-handler.go     # HTTP endpoints (if any)
│   │   └── grpc-handler/
│   │       └── user.go             # gRPC service implementation
│   ├── infrastructure/
│   │   ├── database/
│   │   │   └── postgresql.go       # PostgreSQL connection and queries
│   │   ├── grpc-client/
│   │   │   └── access-client.go    # gRPC client for auth-access (for role assignment)
│   │   ├── hasher/
│   │   │   └── bcrypt.go           # Password hashing (bcrypt/Argon2)
│   │   ├── secrets/
│   │   │   ├── provider.go         # Secret provider interface
│   │   │   └── secrets.go          # Vault-based secret retrieval
│   │   └── vault-client/
│   │       └── vault.go            # Vault client wrapper
│   ├── repository/
│   │   ├── repository.go           # Data access layer (CRUD for users)
│   │   └── tools.go                # Repository helpers
│   ├── server/
│   │   ├── grpc-server/
│   │   │   └── grpc-server.go      # gRPC server setup
│   │   └── http-server/
│   │       └── http-server.go      # HTTP server setup
│   └── service/
│       ├── service.go              # Business logic (user management)
│       ├── service_test.go         # Unit tests
│       └── mocks_test.go           # Mock implementations
├── migrations/
│   ├── 001_create_user_table.up.sql    # Initial schema
│   └── 001_create_user_table.down.sql  # Rollback schema
├── Dockerfile                      # Docker build instructions
├── go.mod
└── go.sum
```

## Configuration

The service uses a YAML configuration file and environment variables. Two presets are provided:

configs/config-local.yml – for local development (HTTP 9082, gRPC 9012).

configs/config-prod.yml – for production (HTTP 8082, gRPC 8012).

You can copy and edit the appropriate file. The configuration path can be passed via the -config flag or the CONFIG_PATH environment variable.

### Environment Variables

These variables **must** be set:

| Variable             | Description                                                 |
| -------------------- | ----------------------------------------------------------- |
| `VAULT_SECRET_TOKEN` | Vault authentication token.                                 |
| `VAULT_ADDRESS`      | Vault server address (e.g., `http://vault:8200`).           |
| `VAULT_MOUNT`        | Vault mount path (e.g., `secret`).                          |
| `USER_POSTGRES_PATH` | ault path for PostgreSQL credentials (e.g., user/postgres). |
| `HASHER_PATH`        | Vault path for hashing parameters (e.g., user/hasher).      |

The service fetches secrets from Vault using these paths. For example, `TOKEN_PATH` must contain a key `jwt_secret` or similar (depending on your implementation).

The service fetches database credentials from Vault using the USER_POSTGRES_PATH. The secret must contain the following keys:

```json
{
	"host": "user_db",
	"port": "5433",
	"username": "admin",
	"password": "securepassword",
	"dbname": "auth_user",
	"sslmode": "disable"
}
```

### Configuration File Details

The YAML configuration includes:

http – HTTP server settings (port, timeouts, base path).

grpc – gRPC server port and connection details for the auth-access service (access_api_host, access_api_port). This is used to assign roles to users.

settings – default timeouts and hash_cost (cost factor for password hashing, e.g., bcrypt cost).

db – PostgreSQL connection details (host, port, ssl_mode). Note that credentials are taken from Vault.

kafka – (optional, only in prod) Kafka brokers and topic for publishing user-role events (e.g., when a user's role changes). This can be used for asynchronous notifications.

Example snippet from config-local.yml:

```yaml
settings:
  default_timeout: 15s
  hash_cost: 12

db:
  host: 'user_db'
  port: 5433
  ssl_mode: disable
```

## Main Responsibilities

### User Management

The service manages the complete lifecycle of a user.

Typical operations include:

- Create user
- Get user
- Update user
- Delete user
- Find user by ID
- Find user by email or username

The exact operations depend on the API exposed by the service.

## Running the Service

### Local Development

Build and run:

```bash
go run cmd/main.go -config=configs/config-local.yml
```

## Dependencies

| Dependency | Purpose                             |
| ---------- | ----------------------------------- |
| PostgreSQL | Persistent user storage             |
| Vault      | Secrets and sensitive configuration |
| gRPC       | Service-to-service communication    |
| Password   | Hasher Secure password storage      |

## Security

The service follows several security principles:s

- Passwords are never stored in plain text.
- Passwords are hashed before persistence.
- Database credentials are not stored in source code.
- Secrets are managed through HashiCorp Vault.
- The PostgreSQL database is owned exclusively by the service.
- Other services access user data through the gRPC API.
- Sensitive user data should only be returned when required by the calling service.

## License

This project is licensed under the MIT License – see the [LICENSE](LICENSE) file for details.

You are free to use, modify, distribute, and sublicense the code for both commercial and non‑commercial purposes, provided that the original copyright notice and permission notice are included in all copies or substantial portions of the software.

For more information, see the full [MIT License](https://opensource.org/licenses/MIT).
