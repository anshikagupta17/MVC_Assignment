CREATE TABLE villages (
    id SERIAL PRIMARY KEY,
    townhall_level INT NOT NULL DEFAULT 1 check (townhall_level>=1 and townhall_level<=4),
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    gold INT NOT NULL DEFAULT 750 check(gold>=0),
    elixir INT NOT NULL DEFAULT 750 check (elixir>=0),
    housing_space INT NOT NULL,
    trophies INT NOT NULL DEFAULT 0,
    layout JSONB not null
);