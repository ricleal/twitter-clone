BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id varchar(36) PRIMARY KEY,
    username varchar(128) UNIQUE NOT NULL,
    email varchar(128) UNIQUE NOT NULL,
    name varchar(256),
    deleted_at timestamptz,
    created_at timestamptz,
    updated_at timestamptz
);

CREATE TABLE IF NOT EXISTS tweets (
    id varchar(36) PRIMARY KEY,
    user_id varchar(36) REFERENCES users(id) NOT NULL,
    content text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz,
    updated_at timestamptz
);

COMMIT;
