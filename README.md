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

## Dependencies

| Dependency | Purpose                             |
| ---------- | ----------------------------------- |
| PostgreSQL | Persistent user storage             |
| Vault      | Secrets and sensitive configuration |
| gRPC       | Service-to-service communication    |
| Password   | Hasher Secure password storage      |

## Security

The service follows several security principles:

- Passwords are never stored in plain text.
- Passwords are hashed before persistence.
- Database credentials are not stored in source code.
- Secrets are managed through HashiCorp Vault.
- The PostgreSQL database is owned exclusively by the service.
- Other services access user data through the gRPC API.
- Sensitive user data should only be returned when required by the calling service.
