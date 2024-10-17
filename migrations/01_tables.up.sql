
-- urls - хранит список уникальных URL и их коротких ключей
CREATE TABLE IF NOT EXISTS urls (
    short_id TEXT PRIMARY KEY,         -- Короткий ключ
    original_url TEXT NOT NULL,        -- Оригинальный URL
    UNIQUE (original_url)
);




