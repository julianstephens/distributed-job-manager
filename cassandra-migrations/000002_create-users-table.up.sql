CREATE TABLE IF NOT EXISTS users(
  id text,
  username text,
  password text,
  PRIMARY KEY (id, username)
);
