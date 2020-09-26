CREATE TABLE player_score_dg_tmp
(
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `name` VARCHAR(255) NOT NULL,
    `score` INTEGER NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX player_name ON player_score_dg_tmp (name);

INSERT INTO player_score_dg_tmp(id, name, score, created_at, updated_at) SELECT id, name, score, created_at, updated_at FROM player_score;

DROP TABLE player_score;

ALTER TABLE player_score_dg_tmp RENAME TO player_score;

