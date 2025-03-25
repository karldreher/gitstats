-- Commmit table contains all commit messages and their dates
CREATE TABLE IF NOT EXISTS commits (
    id SERIAL PRIMARY KEY,
    commit VARCHAR(64) NOT NULL,
    conventional BOOLEAN GENERATED ALWAYS AS (commit ~ '^(feat|fix|docs|style|refactor|perf|test|chore)(\(.+\))?: .{1,50}') STORED,
    repo VARCHAR(255) NOT NULL,
    date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
