CREATE TABLE army_metadata (
    type_id INT NOT NULL REFERENCES buildings_metadata(id) ON DELETE CASCADE,
    level INT NOT NULL CHECK(level>=0 AND level<=4),
    capacity INT NOT NULL CHECK(capacity>=0),
    PRIMARY KEY(type_id, level)
);