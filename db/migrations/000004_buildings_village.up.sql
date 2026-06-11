CREATE TABLE buildings_village (
    id SERIAL PRIMARY KEY,
    village_id INT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    level INT NOT NULL CHECK (level>=1 AND level<=4),
    building_id int not null references buildings_metadata(id) on delete cascade,
    upgrade_ends_at timestamp null,
    x int not null default -1,
    y int not null default -1
);