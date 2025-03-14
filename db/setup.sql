-- Commmit table contains all commit messages and their dates
CREATE TABLE IF NOT EXISTS commits (
    id SERIAL PRIMARY KEY,
    commit VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
