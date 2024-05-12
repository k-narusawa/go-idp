CREATE TABLE refresh_tokens (
  signature VARCHAR(255) NOT NULL UNIQUE,
  request_id VARCHAR(40) NOT NULL,
  client_id VARCHAR(255) NOT NULL,
  requested_at TIMESTAMP NOT NULL,
  scope VARCHAR(255) NOT NULL,
  granted_scope VARCHAR(255) NOT NULL,
  form_data TEXT NOT NULL,
  session_data TEXT NOT NULL,
  subject VARCHAR(40) NOT NULL,
  active BOOLEAN NOT NULL,
  requested_audience VARCHAR(255) NOT NULL,
  granted_audience VARCHAR(255) NOT NULL,
  PRIMARY KEY (signature),
  FOREIGN KEY (client_id) REFERENCES clients (id)
);
