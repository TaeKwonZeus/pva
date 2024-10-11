PRAGMA foreign_keys= ON;

CREATE TABLE IF NOT EXISTS users
(
    id                    INTEGER PRIMARY KEY,
    username              TEXT NOT NULL UNIQUE,
    role                  TEXT NOT NULL,

    salt                  TEXT NOT NULL,
    public_key            TEXT NOT NULL,
    private_key_encrypted TEXT NOT NULL
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
    name               TEXT    NOT NULL,
    description        TEXT    NOT NULL,
    password_encrypted TEXT    NOT NULL,
    created_at         INTEGER NOT NULL,
    updated_at         INTEGER NOT NULL,
    vault_id           INTEGER REFERENCES vaults (id) ON DELETE CASCADE,

    UNIQUE (name, vault_id)
);

CREATE TABLE IF NOT EXISTS vault_keys
(
    user_id             INTEGER REFERENCES users (id) ON DELETE CASCADE,
    vault_id            INTEGER REFERENCES vaults (id) ON DELETE CASCADE,

    -- Encrypted with user's public key
    vault_key_encrypted TEXT NOT NULL,

    PRIMARY KEY (user_id, vault_id)
);

CREATE TABLE IF NOT EXISTS devices
(
    id          INTEGER PRIMARY KEY,
    ip          TEXT UNIQUE NOT NULL,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS documents
(
    id                INTEGER PRIMARY KEY,
    name              TEXT UNIQUE NOT NULL,
    payload_encrypted TEXT        NOT NULL,
    created_at        INTEGER     NOT NULL,
    updated_at        INTEGER     NOT NULL
);

CREATE TABLE IF NOT EXISTS attachments
(
    id          INTEGER PRIMARY KEY,
    document_id INTEGER REFERENCES documents (id) ON DELETE CASCADE,
    name        TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS document_keys
(
    user_id                INTEGER REFERENCES users (id) ON DELETE CASCADE,
    document_id            INTEGER REFERENCES documents (id) ON DELETE CASCADE,

    document_key_encrypted TEXT NOT NULL,

    PRIMARY KEY (user_id, document_id)
);
