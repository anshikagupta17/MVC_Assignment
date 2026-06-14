ALTER TABLE battles
ADD COLUMN destruction_percent INT NOT NULL CHECK(destruction_percent >= 0 AND destruction_percent <= 100),
ADD COLUMN trophy_change INT NOT NULL;