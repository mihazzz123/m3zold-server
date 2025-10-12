-- Создание схемы
CREATE SCHEMA IF NOT EXISTS m3zold_schema;

-- Включение расширения для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создание таблицы для отслеживания выполненных миграций
CREATE TABLE IF NOT EXISTS m3zold_schema.schema_migrations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    executed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индекс для быстрого поиска миграций по имени
CREATE INDEX IF NOT EXISTS idx_schema_migrations_name ON m3zold_schema.schema_migrations(name);