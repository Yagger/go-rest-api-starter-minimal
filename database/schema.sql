CREATE TABLE IF NOT EXISTS account (
  account_id VARCHAR(36) PRIMARY KEY,
  admin_id VARCHAR(36) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  salt VARCHAR(36) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL UNIQUE,
  last_login DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE CURRENT_TIMESTAMP,
  deleted_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
  full_name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  FOREIGN KEY (admin_id) REFERENCES account(account_id)
);

CREATE TABLE IF NOT EXISTS role (
  role VARCHAR(16) PRIMARY KEY
);

INSERT IGNORE INTO role (role) VALUES ('user'), ('admin'), ('superadmin');

CREATE TABLE IF NOT EXISTS account_role (
  account_id VARCHAR(36) NOT NULL,
  role VARCHAR(16) NOT NULL,
  PRIMARY KEY (account_id, role),
  FOREIGN KEY (account_id) REFERENCES account(account_id) ON DELETE CASCADE,
  FOREIGN KEY (role) REFERENCES role(role) ON UPDATE CASCADE
);

CREATE VIEW account_view AS
SELECT account_id, admin_id, email, last_login, created_at, updated_at, deleted_at, full_name, GROUP_CONCAT(role) AS roles
FROM account JOIN account_role USING (account_id)
GROUP BY created_at;