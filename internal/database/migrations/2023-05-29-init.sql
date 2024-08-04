CREATE TABLE IF NOT EXISTS credentials (
    id uuid PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz DEFAULT NULL
);
