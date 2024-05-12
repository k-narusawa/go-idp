CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  user_id VARCHAR(36) NOT NULL UNIQUE,
  username VARCHAR(255) NOT NULL UNIQUE,
  password BYTEA
);

CREATE INDEX idx_user_id ON users (user_id);
CREATE INDEX idx_username ON users (username);
