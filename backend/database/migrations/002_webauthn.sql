CREATE TABLE webauthn_sessions (
                                   id SERIAL PRIMARY KEY,
                                   email TEXT NOT NULL,

                                   challenge BYTEA NOT NULL,
                                   data JSONB NOT NULL,

                                   expires_at TIMESTAMP NOT NULL
);