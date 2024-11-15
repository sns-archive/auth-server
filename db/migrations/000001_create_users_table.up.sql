CREATE TABLE users (
  id BINARY(16) PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  -- 後で暗号化して保存する
  password VARCHAR(255) NOT NULL
);
