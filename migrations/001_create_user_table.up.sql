CREATE TABLE auth_user
(
    id varchar(255) PRIMARY KEY NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,
    username varchar(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    fullname text NOT NULL,
    is_active boolean NOT NULL DEFAULT TRUE,
    is_superuser boolean NOT NULL DEFAULT FALSE,
    is_banned boolean NOT NULL DEFAULT FALSE,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);