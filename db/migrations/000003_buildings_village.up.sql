CREATE TABLE buildings_village (
    id SERIAL PRIMARY KEY,
    village_id INT not null references villages(id) on delete cascade,
    level int not null check (level>=1 and level<=4),
    building_id int not null references buildings_metadata(id) on delete cascade,
    upgrade_ends_at timestamp null,
    x int not null default -1,
    y int not null default -1
);