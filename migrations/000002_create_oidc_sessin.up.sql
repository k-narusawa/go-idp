CREATE TABLE oidc_sessions (
  signature varchar(255) NOT NULL UNIQUE,
  request_id varchar(40) NOT NULL,
  client_id varchar(255) NOT NULL,
  requested_at timestamp NOT NULL,
  scope varchar(255) NOT NULL,
  granted_scope varchar(255) NOT NULL,
  form_data text NOT NULL,
  session_data text NOT NULL,
  subject text NOT NULL,
  active boolean NOT NULL,
  requested_audience varchar(255) NOT NULL,
  granted_audience varchar(255) NOT NULL,
  PRIMARY KEY (signature),
  FOREIGN KEY (client_id) REFERENCES clients(id)
);