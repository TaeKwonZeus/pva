CREATE TABLE IF NOT EXISTS users
(
    id                    INTEGER PRIMARY KEY,
    username              TEXT NOT NULL UNIQUE,
    salt                  TEXT NOT NULL,
    public_key            TEXT NOT NULL,
    private_key_encrypted TEXT NOT NULL,
    role                  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS vaults
(
    id       INTEGER PRIMARY KEY,
    name     TEXT    NOT NULL UNIQUE,
    owner_id INTEGER REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS passwords
(
    id                 INTEGER PRIMARY KEY,
    name               TEXT                              NOT NULL,
    description        TEXT,
    password_encrypted TEXT                              NOT NULL,
    created_at         NUMERIC DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at         NUMERIC DEFAULT CURRENT_TIMESTAMP NOT NULL,
    vault_id           INTEGER REFERENCES vaults (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS vault_access
(
    user_id  INTEGER REFERENCES users (id) ON DELETE CASCADE,
    vault_id INTEGER REFERENCES vaults (id) ON DELETE CASCADE
);