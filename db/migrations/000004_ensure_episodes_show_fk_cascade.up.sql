ALTER TABLE episodes
DROP CONSTRAINT IF EXISTS episodes_show_id_fkey;

ALTER TABLE episodes
ADD CONSTRAINT episodes_show_id_fkey
FOREIGN KEY (show_id)
REFERENCES shows(internal_show_id)
ON DELETE CASCADE;
