-- For each new user create a keypair
-- Store public key here
-- Store private key on user's device encrypted with password (AES)
CREATE TABLE IF NOT EXISTS users
(
    id         INTEGER PRIMARY KEY,
    username   TEXT UNIQUE NOT NULL,
    public_key TEXT        NOT NULL
);

CREATE TABLE IF NOT EXISTS permissions
(
    user_id    INTEGER,
    permission TEXT,
    PRIMARY KEY (user_id, permission),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- For each new secret create a random key
-- We don't store the key: we encrypt the secret with it and
-- for every user with read permission create a row in secret_keys
-- Store the key there encrypted with each user's public key
CREATE TABLE IF NOT EXISTS secrets
(
    id               INTEGER PRIMARY KEY,
    -- Encrypted with randomly generated key
    encrypted_secret BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS secret_keys
(
    secret_id     INTEGER,
    user_id       INTEGER,
    -- Random key encrypted with the user's public key
    encrypted_key BLOB NOT NULL,
    primary key (secret_id, user_id),
    foreign key (secret_id) references secrets (id),
    foreign key (user_id) references users (id)
);
