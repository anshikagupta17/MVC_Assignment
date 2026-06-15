ALTER TABLE mines_metadata
ALTER COLUMN production_rate TYPE INT USING production_rate::INT;