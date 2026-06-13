CREATE TABLE building_limits (
    building_id INT NOT NULL REFERENCES buildings_metadata(id),
    townhall_level INT NOT NULL,
    max_quantity INT NOT NULL,
    PRIMARY KEY(building_id, townhall_level)
);