CREATE TABLE battles(
    id SERIAL PRIMARY KEY,
    attacker_id INT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    defender_id INT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    loot_gold INT NOT NULL,
    loot_elixir INT NOT NULL,
    start_time TIMESTAMP,
    stars INT NOT NULL CHECK (stars>=0 AND stars<=3)
);