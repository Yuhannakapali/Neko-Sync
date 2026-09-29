-- Runs once when the Postgres container initialises an empty data volume.
-- The schema lives in versioned migrations under migrations/; apply it with
-- `make migrate-up` (or the `migrate` compose service).
SELECT 1;
