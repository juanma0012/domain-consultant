CREATE USER IF NOT EXISTS maxroach;
CREATE DATABASE IF NOT EXISTS logs;
GRANT ALL ON DATABASE logs TO maxroach;
USE logs;
 DROP TABLE  Server;
 DROP TABLE  History;
 DROP TABLE  Response;
CREATE TABLE IF NOT EXISTS Response
(
  response_id SERIAL PRIMARY KEY,
  domain VARCHAR NOT NULL,
  servers_changed BOOL,
  ssl_grade VARCHAR(2),
  previous_ssl_grade VARCHAR(2),
  logo VARCHAR,
  title VARCHAR,
  is_down BOOL,
  user_session_id VARCHAR NOT NULL,
  created VARCHAR NOT NULL,
);

CREATE TABLE IF NOT EXISTS Server
(
  server_id SERIAL PRIMARY KEY,
  address VARCHAR NOT NULL,
  ssl_grade VARCHAR(2),
  country VARCHAR,
  owner VARCHAR,
  response_id INT NOT NULL,
  FOREIGN KEY (response_id) REFERENCES Response(response_id)
);

CREATE TABLE IF NOT EXISTS History
(
  user_session_id VARCHAR NOT NULL,
  created VARCHAR NOT NULL,
  response_id INT NOT NULL,
  FOREIGN KEY (response_id) REFERENCES Response(response_id)
);