DROP TABLE IF EXISTS game_game_option_value;
DROP TABLE IF EXISTS game_type_game_option;
DROP TABLE IF EXISTS game_option;
ALTER TABLE game REMOVE COLUMN game_type_id;
DROP TABLE IF EXISTS game_type;
