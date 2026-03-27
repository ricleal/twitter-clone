BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuidv7() PRIMARY KEY,
    username varchar(128) UNIQUE NOT NULL,
    email varchar(128) UNIQUE NOT NULL,
    name varchar(256),
    deleted_at timestamptz,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS tweets (
    id uuid DEFAULT uuidv7() PRIMARY KEY,
    user_id uuid REFERENCES users(id) NOT NULL,
    content text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz,
    updated_at timestamptz
);

-- tweets by user (e.g. user timeline)
CREATE INDEX IF NOT EXISTS idx_tweets_user_id ON tweets (user_id);

-- timeline ordering / pagination
CREATE INDEX IF NOT EXISTS idx_tweets_created_at ON tweets (created_at);

-- soft-delete filtering
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_tweets_deleted_at ON tweets (deleted_at) WHERE deleted_at IS NULL;

COMMIT;
