-- Instance schema v1: libraries and media files.
-- Times are INTEGER unix nanoseconds so mod_time compares exactly with the filesystem.
-- Composite fields (streams, provider_ids) are JSON text.

CREATE TABLE libraries (
    id         TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    kind       TEXT NOT NULL CHECK (kind IN ('movie', 'series', 'anime')),
    root       TEXT NOT NULL UNIQUE,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

CREATE TABLE media_files (
    id           TEXT PRIMARY KEY,
    library_id   TEXT NOT NULL REFERENCES libraries (id) ON DELETE CASCADE,
    path         TEXT NOT NULL UNIQUE,
    size         INTEGER NOT NULL,
    mod_time     INTEGER NOT NULL,
    kind         TEXT NOT NULL CHECK (kind IN ('movie', 'series', 'anime')),
    container    TEXT NOT NULL DEFAULT '',
    duration     REAL NOT NULL DEFAULT 0,
    streams      TEXT NOT NULL DEFAULT '[]',
    parsed_title TEXT NOT NULL,
    parsed_year  INTEGER,
    season       INTEGER,
    episode      REAL,
    provider_ids TEXT NOT NULL DEFAULT '{}',
    match_status TEXT NOT NULL DEFAULT 'unmatched'
                 CHECK (match_status IN ('unmatched', 'matched', 'manual')),
    created_at   INTEGER NOT NULL,
    updated_at   INTEGER NOT NULL
);

CREATE INDEX media_files_library_id   ON media_files (library_id);
CREATE INDEX media_files_match_status ON media_files (match_status);
