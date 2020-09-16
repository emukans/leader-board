CREATE TABLE player_score (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(255) UNIQUE NOT NULL,
    `score` INTEGER NOT NULL,
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);


CREATE TRIGGER tg_player_score_updated_at
AFTER UPDATE
ON player_score FOR EACH ROW
BEGIN
  UPDATE player_score SET updated_at = current_timestamp
    WHERE id = old.id;
END;
