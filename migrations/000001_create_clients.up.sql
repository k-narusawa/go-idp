CREATE TABLE clients (
  id VARCHAR(255) NOT NULL UNIQUE,
  secret BYTEA,
  redirect_uris TEXT,
  grant_types TEXT,
  response_types TEXT,
  scopes TEXT,
  audience TEXT,
  public BOOLEAN NOT NULL,
  PRIMARY KEY (id)
);
