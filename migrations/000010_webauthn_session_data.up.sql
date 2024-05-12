CREATE TABLE webauthn_session_data (
  session_id SERIAL PRIMARY KEY,
  challenge VARCHAR(255) UNIQUE NOT NULL,
  user_id BYTEA,
  allowed_credential_ids TEXT,
  expires TIMESTAMP NOT NULL,
  user_verification TEXT NOT NULL,
  extensions TEXT
);

CREATE INDEX idx_webauthn_session_data_challenge ON webauthn_session_data (challenge);
