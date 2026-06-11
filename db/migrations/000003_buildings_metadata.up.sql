CREATE TYPE currency AS ENUM (
    'Gold',
    'Elixir'
);

CREATE TABLE buildings_metadata (
    id SERIAL PRIMARY KEY,
    size_x INT NOT NULL CHECK(size_x>=0),
    size_y INT NOT NULL CHECK(size_y>=0),
    name VARCHAR(30) UNIQUE NOT NULL,
    upgrade_cost INT NOT NULL,
    cost_type currency NOT NULL,
    upgrade_time_sec INT NOT NULL
);