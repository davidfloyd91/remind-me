CREATE TABLE users(
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) NOT NULL,
  email VARCHAR(200),
  password char(60),
  created_at TIMESTAMPTZ DEFAULT current_timestamp,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

CREATE TABLE events(
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users (id),
  name VARCHAR(200) NOT NULL,
  description VARCHAR(500),
  scheduled TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ DEFAULT current_timestamp,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

CREATE OR REPLACE FUNCTION updated_at_column() RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = current_timestamp;
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER users_updated_at_column
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE PROCEDURE updated_at_column();
