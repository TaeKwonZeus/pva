-- For each new user create a keypair
-- Store public key here
-- Store private key encrypted with password (AES)
CREATE TABLE IF NOT EXISTS users
(
    id                    INTEGER PRIMARY KEY,
    username              TEXT UNIQUE NOT NULL,
    public_key            TEXT        NOT NULL,
    encrypted_private_key TEXT        NOT NULL
);

-- For each new password create a random key
-- We don't store the key: we encrypt the password with it and
-- for every user with read permission create a row in password_keys
-- Store the key there encrypted with each user's public key
CREATE TABLE IF NOT EXISTS secrets
(
    id               INTEGER PRIMARY KEY,
    -- Encrypted with randomly generated key
    encrypted_secret TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS secret_keys
(
    secret_id     INTEGER,
    user_id       INTEGER,
    -- Random key encrypted with the user's public key
    encrypted_key TEXT NOT NULL,
    primary key (secret_id, user_id),
    foreign key (secret_id) references secrets (id),
    foreign key (user_id) references users (id)
);
