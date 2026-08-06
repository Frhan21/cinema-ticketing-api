ALTER TABLE movies
ADD COLUMN tmdb_id BIGINT;

CREATE UNIQUE INDEX movies_tmdb_id_unique
ON movies (tmdb_id)
WHERE tmdb_id IS NOT NULL;
