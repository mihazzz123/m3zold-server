-- Создание таблицы устройств
CREATE TABLE IF NOT EXISTS m3zold_schema.devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES m3zold_schema.users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT,
    status TEXT DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Создание индексов для таблицы devices
CREATE INDEX IF NOT EXISTS idx_devices_user_id ON m3zold_schema.devices(user_id);
CREATE INDEX IF NOT EXISTS idx_devices_status ON m3zold_schema.devices(status);
CREATE INDEX IF NOT EXISTS idx_devices_type ON m3zold_schema.devices(type);