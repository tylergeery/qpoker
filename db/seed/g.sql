INSERT INTO game_type
    (id, key, display_name)
VALUES
    (1, 'holdem', 'Texas Holdem (Poker)'),
    (2, 'hearts', 'Hearts');


INSERT INTO game_option
    (id, name, label, type)
VALUES
    (1, 'capacity', 'Capacity', 'integer'),
    (2, 'time_between_hands', 'Time Between Hands (s)', 'integer'),
    (3, 'decision_time', 'Player Decision Time (s)', 'integer'),
    (4, 'big_blind', 'Big Blind', 'integer'),
    (5, 'buy_in_min', 'Min Buy In', 'integer'),
    (6, 'buy_in_max', 'Max Buy In', 'integer');


INSERT INTO game_type_game_option
    (id, game_type_id, game_option_id, default_value)
VALUES
    (1, 1, 1, '12'),
    (2, 1, 2, '5'),
    (3, 1, 3, '30'),
    (4, 1, 4, '50'),
    (5, 1, 5, '500'),
    (6, 1, 6, '5000'),
    (7, 2, 1, '4'),
    (8, 2, 2, '5'),
    (9, 2, 3, '30');
