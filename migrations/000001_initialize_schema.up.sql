CREATE TABLE proxy_profiles
(
    id      INTEGER PRIMARY KEY,
    name    TEXT    NOT NULL UNIQUE,
    type    INTEGER NOT NULL CHECK (type >= 1 AND type <= 4),
    address TEXT
);

CREATE TABLE rules
(
    id               INTEGER PRIMARY KEY,
    regex            TEXT                                   NOT NULL,
    proxy_profile_id INTEGER REFERENCES proxy_profiles (id) NOT NULL
);
