CREATE TABLE login_skip_sessions (
  session_id BIGSERIAL PRIMARY KEY,
  token VARCHAR(255) NOT NULL UNIQUE,
  user_id VARCHAR(255),
  expires_at TIMESTAMP,
  CONSTRAINT token_index UNIQUE (token)
);

CREATE INDEX idx_token ON login_skip_sessions (token);
