ALTER TABLE buildings_village
ADD COLUMN last_collected_at TIMESTAMP NULL;

UPDATE buildings_village
SET last_collected_at = NOW()
WHERE building_id IN (6,7);