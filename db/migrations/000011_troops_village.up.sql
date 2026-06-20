CREATE TABLE troops_village(
    village_id INT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    troops_id INT NOT NULL REFERENCES troops_base_metadata(id) ON DELETE CASCADE,
    level INT NOT NULL CHECK (level>=0 AND level<=4),
    quantity INT NOT NULL CHECK(quantity>=0),
    upgrade_ends_at timestamp,
    PRIMARY KEY(village_id, troops_id)
);