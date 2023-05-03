CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
    id uuid PRIMARY KEY,
    email varchar(255) NOT NULL,
    is_active boolean NOT NULL,
    refresh_token varchar NOT NULL UNIQUE,
    user_agent varchar(255) NOT NULL,
    client_ip varchar(255) NOT NULL,
    created_at timestamptz NOT NULL,
    expires_at timestamptz NOT NULL
);

ALTER TABLE "sessions" ADD FOREIGN KEY (email) REFERENCES users(email) ON DELETE CASCADE;

