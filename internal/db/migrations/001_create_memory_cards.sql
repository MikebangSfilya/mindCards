CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY, 
    email VARCHAR(100) NOT NULL UNIQUE, 
    password_hash VARCHAR(250) NOT NULL
);


CREATE TABLE IF NOT EXISTS memory_cards (
    card_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    card_description TEXT,
    tag VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    level_study SMALLINT DEFAULT 0 CHECK (level_study >= 0 AND level_study <= 5),
    learned BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);



CREATE INDEX IF NOT EXISTS idx_tag ON memory_cards (tag);
CREATE INDEX IF NOT EXISTS idx_title ON memory_cards (title);
CREATE INDEX IF NOT EXISTS idx_learned ON memory_cards (learned);
CREATE INDEX IF NOT EXISTS idx_level_study ON memory_cards (level_study);
CREATE INDEX IF NOT EXISTS idx_memory_cards_user_id ON memory_cards (user_id);
CREATE INDEX IF NOT EXISTS idx_user_learned ON memory_cards (user_id, learned);
CREATE INDEX IF NOT EXISTS idx_user_tag ON memory_cards (user_id, tag);