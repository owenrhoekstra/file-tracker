CREATE TABLE users (
                       id BYTEA PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL
);

CREATE TABLE credentials (
                             id SERIAL PRIMARY KEY,
                             user_id BYTEA NOT NULL REFERENCES users(id) ON DELETE CASCADE,

                             credential_id BYTEA NOT NULL UNIQUE,
                             public_key BYTEA NOT NULL,
                             attestation_type TEXT NOT NULL,
                             transports TEXT[],

                             sign_count BIGINT NOT NULL DEFAULT 0
);