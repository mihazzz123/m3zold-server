-- Создание таблицы токенов аутентификации
CREATE TABLE IF NOT EXISTS m3zold_schema.auth_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES m3zold_schema.users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    token_type TEXT NOT NULL CHECK (token_type IN ('access', 'refresh')),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    blacklisted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для таблицы auth_tokens
CREATE INDEX IF NOT EXISTS idx_auth_tokens_token ON m3zold_schema.auth_tokens(token);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_user_id ON m3zold_schema.auth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_expires_at ON m3zold_schema.auth_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_auth_tokens_blacklisted ON m3zold_schema.auth_tokens(blacklisted);