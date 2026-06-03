CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    pass_hash VARCHAR(255) NOT NULL, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE villages (
    id SERIAL PRIMARY KEY,
    townhall_level INT NOT NULL DEFAULT 1 check (townhall_level>=1 and townhall_level<=4),
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    gold BIGINT NOT NULL DEFAULT 750 check(gold>=0),
    elixir BIGINT NOT NULL DEFAULT 750 check (elixir>=0),
    housing_space INT NOT NULL
);

CREATE TABLE troops_metadata (
    id serial primary key,
    type_id int NOT NULL,
    level int NOT NULL, 
    damage INT NOT NULL,
    max_health INT NOT NULL Check(max_health>=0),
    name VARCHAR(20) NOT NULL,
    upgrade_cost BIGINT NOT NULL,
    housing_space INT NOT NULL,
    cost_type VARCHAR(10) NOT NULL DEFAULT 'Elixir',
    speed INT NOT NULL,
    unlock_level int not null,
    UNIQUE(TYPE_ID,level)
);

CREATE TABLE troops_village (
    id SERIAL PRIMARY KEY,
    village_id BIGINT  NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    quantity BIGINT DEFAULT 0 check(quantity>=0),
    troops_id BIGINT NOT NULL REFERENCES troops_metadata(id) ON DELETE CASCADE ,
    UNIQUE(village_id, troops_id)
);

CREATE TABLE defense_metadata (
    id serial primary key,
    type_id int NOT NULL,
    level int NOT NULL,
    name VARCHAR(20) NOT NULL,
    unlock_level INT NOT NULL,
    damage BIGINT NOT NULL,
    max_health INT NOT NULL CHECK (max_health>=0),
    upgrade_cost bigint not null,
    range INT NOT NULL,
    cost_type VARCHAR(10) DEFAULT 'Gold' NOT NULL,
    max_quantity BIGINT NOT NULL,
    UNIQUE(type_id,level)
);

CREATE TABLE defense_village (
    id SERIAL PRIMARY KEY,
    village_id BIGINT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    X BIGINT NOT NULL,
    Y BIGINT NOT NULL,
    defense_id BIGINT NOT NULL REFERENCES defense_metadata(id) ON DELETE CASCADE,
    UNIQUE(village_id, X, Y)
);

CREATE TABLE battles (
    ID SERIAL PRIMARY KEY,
    attacker_id BIGINT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    defender_id BIGINT NOT NULL REFERENCES villages(id) ON DELETE CASCADE,
    loot_gold BIGINT NOT NULL check(loot_gold>=0),
    loot_elixir BIGINT NOT NULL CHECK(loot_elixir>=0),
    start_time TIMESTAMP NOT NULL,
    stars INT NOT NULL CHECK(stars>=0 AND stars<=3)
);

CREATE TABLE battle_events (
    battle_id BIGINT NOT NULL references battles(id) ON DELETE CASCADE,
    event_order BIGINT NOT NULL, 
    event_text JSONB NOT NULL,
    primary key(battle_id, event_order)
);