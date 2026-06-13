CREATE TABLE troops_base_metadata(
    id SERIAL PRIMARY KEY,
    name VARCHAR(20) NOT NULL, 
    cost_type currency NOT NULL DEFAULT 'Elixir',
    speed INT NOT NULL, 
    unlock_level INT NOT NULL CHECK(unlock_level>=0 and unlock_level<=4),
    upgrade_time_sec INT NOT NULL,
    housing_space INT NOT NULL
);