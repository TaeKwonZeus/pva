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

-- For each new password create a random key
-- We don't store the key: we encrypt the password with it and
-- for every user with read permission create a row in password_keys
-- Store the key there encrypted with each user's public key
CREATE TABLE IF NOT EXISTS passwords
(
    id               INTEGER PRIMARY KEY,
    -- Encrypted with randomly generated key
    encrypted_password BLOB NOT NULL
);

CREATE TABLE IF NOT EXISTS password_keys
(
    password_id     INTEGER,
    user_id       INTEGER,
    -- Random key encrypted with the user's public key
    encrypted_key BLOB NOT NULL,
    primary key (password_id, user_id),
    foreign key (password_id) references passwords (id),
    foreign key (user_id) references users (id)
);
