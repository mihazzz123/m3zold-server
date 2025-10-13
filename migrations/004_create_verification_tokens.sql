-- Создание таблицы токенов верификации
CREATE TABLE IF NOT EXISTS m3zold_schema.verification_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES m3zold_schema.users(id) ON DELETE CASCADE,
    token TEXT UNIQUE NOT NULL,
    token_type TEXT DEFAULT 'email_verification',
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Индексы для таблицы verification_tokens
CREATE INDEX IF NOT EXISTS idx_verification_tokens_token ON m3zold_schema.verification_tokens(token);
CREATE INDEX IF NOT EXISTS idx_verification_tokens_user_id ON m3zold_schema.verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_verification_tokens_expires_at ON m3zold_schema.verification_tokens(expires_at);