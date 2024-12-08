
-- urls - хранит список уникальных URL и их коротких ключей
CREATE TABLE IF NOT EXISTS urls (
    idx INTEGER,                      -- Индекс записи в memstore
    short_id TEXT PRIMARY KEY,         -- Короткий ключ
    original_url TEXT NOT NULL,        -- Оригинальный URL
    user_id TEXT,                      -- Идентификатор пользователя
    deleted INTEGER DEFAULT 0,         -- Флаг удаления
    UNIQUE (user_id, original_url)
);


CREATE TABLE IF NOT EXISTS urls1 (
    idx INTEGER,                      -- Индекс записи в memstore
    short_id TEXT PRIMARY KEY,         -- Короткий ключ
    original_url TEXT NOT NULL,        -- Оригинальный URL
    user_id TEXT,                      -- Идентификатор пользователя
    deleted INTEGER DEFAULT 0,         -- Флаг удаления
    UNIQUE (user_id, original_url)
);




