CREATE TABLE mines_metadata(
    type_id INT NOT NULL REFERENCES buildings_metadata(id) ON DELETE CASCADE,
    level INT NOT NULL CHECK(level>=0 AND level<=4),
    production_rate INT NOT NULL CHECK (production_rate>=0),
    PRIMARY KEY(type_id, level)
);