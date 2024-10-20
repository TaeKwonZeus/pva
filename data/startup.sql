PRAGMA foreign_keys= ON;

CREATE TABLE IF NOT EXISTS users
(
    id                    INTEGER PRIMARY KEY,
    username              TEXT NOT NULL UNIQUE,
    role                  TEXT NOT NULL,

    salt                  BLOB NOT NULL,
    public_key            BLOB NOT NULL,
    private_key_encrypted BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS vaults
(
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS passwords
(
    id                 INTEGER PRIMARY KEY,
    name               TEXT NOT NULL,
    description        TEXT NOT NULL,
    password_encrypted BLOB NOT NULL,
    vault_id           INTEGER REFERENCES vaults (id) ON DELETE CASCADE,

    UNIQUE (name, vault_id)
);

CREATE TABLE IF NOT EXISTS vault_keys
(
    user_id       INTEGER REFERENCES users (id) ON DELETE CASCADE,
    vault_id      INTEGER REFERENCES vaults (id) ON DELETE CASCADE,

    -- Encrypted with user's public key
    key_encrypted BLOB NOT NULL,

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
    name              TEXT        NOT NULL,
    file_name         TEXT UNIQUE NOT NULL,
    payload_encrypted BLOB        NOT NULL
);

CREATE TABLE IF NOT EXISTS attachments
(
    id          INTEGER PRIMARY KEY,
    document_id INTEGER REFERENCES documents (id) ON DELETE CASCADE,
    name        TEXT        NOT NULL,
    file_name   TEXT UNIQUE NOT NULL,

    UNIQUE (document_id, name)
);

CREATE TABLE IF NOT EXISTS document_keys
(
    user_id       INTEGER REFERENCES users (id) ON DELETE CASCADE,
    document_id   INTEGER REFERENCES documents (id) ON DELETE CASCADE,

    key_encrypted BLOB NOT NULL,

    PRIMARY KEY (user_id, document_id)
);
