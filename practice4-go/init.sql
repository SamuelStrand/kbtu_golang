-- This file is executed automatically by the postgres container on FIRST start
-- (only when the volume is empty).

CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL UNIQUE,
    age int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id serial PRIMARY KEY,
    user_id int REFERENCES users(id),
    action varchar(255) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
