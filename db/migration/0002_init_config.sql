CREATE TABLE config (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(50) UNIQUE NOT NULL,
    `value` VARCHAR(255) UNIQUE NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TRIGGER tg_config_updated_at
AFTER UPDATE
ON config FOR EACH ROW
BEGIN
  UPDATE config SET updated_at = current_timestamp
    WHERE id = old.id;
END;


INSERT INTO config (name, value) VALUES ('auth_token', '7191ba6933d2ea4d775dd31f6ea351abf794b9fe'), ('page_limit', '10');