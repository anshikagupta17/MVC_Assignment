CREATE TABLE troops_level_metadata(
    type_id INT NOT NULL REFERENCES troops_base_metadata(id) ON DELETE CASCADE,
    level INT NOT NULL CHECK(level>=0 and level<=4),
    damage INT NOT NULL CHECK (damage>=0),
    max_health INT NOT NULL CHECK (max_health>=0),
    upgrade_cost INT NOT NULL, 
    PRIMARY KEY(type_id, level)
);