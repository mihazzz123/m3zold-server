-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS "m3zold_schema"."users" (
    id UUID PRIMARY KEY NOT NULL,
    email TEXT UNIQUE NOT NULL,
    user_name TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    first_name TEXT,
    second_name TEXT,
    last_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL
);

-- Создание индексов для таблицы users
CREATE INDEX IF NOT EXISTS idx_users_email ON "m3zold_schema"."users"(email);
CREATE INDEX IF NOT EXISTS idx_users_user_name ON "m3zold_schema"."users"(user_name);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON "m3zold_schema"."users"(is_active);
CREATE INDEX IF NOT EXISTS idx_users_is_verified ON "m3zold_schema"."users"(is_verified);