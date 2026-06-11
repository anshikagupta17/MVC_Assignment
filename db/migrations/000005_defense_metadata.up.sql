CREATE TABLE defense_metadata (
    type_id INT NOT NULL REFERENCES buildings_metadata(id) ON DELETE CASCADE,
    level INT NOT NULL CHECK(level>=0 AND level<=4),
    unlock_level INT NOT NULL CHECK(unlock_level>=0 and unlock_level<=4),
    damage INT NOT NULL CHECK (damage>=0),
    max_health INT NOT NULL CHECK(max_health>=0),
    range INT NOT NULL CHECK(range>=0),
    PRIMARY KEY(type_id, level)
);