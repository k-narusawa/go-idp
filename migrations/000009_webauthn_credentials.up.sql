CREATE TABLE webauthn_credentials (
  credential_id BIGSERIAL PRIMARY KEY,
  id BYTEA,
  user_id VARCHAR(36) NOT NULL,
  public_key BYTEA,
  attestation_type TEXT,
  transport TEXT,
  flags TEXT,
  authenticator TEXT
);

CREATE INDEX idx_webauthn_credentials_id ON webauthn_credentials (id);
CREATE INDEX idx_webauthn_credentials_user_id ON webauthn_credentials (user_id);
