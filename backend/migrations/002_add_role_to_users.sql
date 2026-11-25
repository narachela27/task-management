USE `task-managemnt`;
renderUsersList
ALTER TABLE users
ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'USER' AFTER password,
ADD CONSTRAINT chk_users_role CHECK (role IN ('USER', 'ADMIN')),
ADD INDEX idx_users_role (role);