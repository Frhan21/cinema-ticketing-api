DROP INDEX IF EXISTS movies_tmdb_id_unique;

ALTER TABLE movies
DROP COLUMN IF EXISTS tmdb_id;
