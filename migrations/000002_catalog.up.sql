-- Hub metadata catalog (Work, WorkChild) and the per-user ContentReference
-- registry. Metadata and pointers only: no column here may hold media bytes or
-- a URL to them.

CREATE TABLE IF NOT EXISTS works (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kind TEXT NOT NULL CHECK (kind IN ('anime', 'manga', 'movie', 'series', 'music', 'book')),
    title TEXT NOT NULL,
    titles JSONB NOT NULL DEFAULT '[]',
    year INT,
    synopsis TEXT,
    artwork JSONB NOT NULL DEFAULT '[]',
    genres TEXT[] NOT NULL DEFAULT '{}',
    tags TEXT[] NOT NULL DEFAULT '{}',
    provider_ids JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_works_kind ON works(kind);
CREATE INDEX IF NOT EXISTS idx_works_provider_ids ON works USING GIN (provider_ids jsonb_path_ops);

CREATE TABLE IF NOT EXISTS work_children (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES work_children(id) ON DELETE CASCADE,
    kind TEXT NOT NULL CHECK (kind IN ('season', 'episode', 'volume', 'chapter', 'track')),
    ordinal DOUBLE PRECISION NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    synopsis TEXT,
    duration INT,
    release_date TIMESTAMP WITH TIME ZONE,
    provider_ids JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- Parent-scoped: season 1 and season 2 can both hold "episode 1".
    UNIQUE NULLS NOT DISTINCT (work_id, parent_id, kind, ordinal)
);

CREATE INDEX IF NOT EXISTS idx_work_children_work_id ON work_children(work_id);

CREATE TABLE IF NOT EXISTS content_references (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    work_id UUID NOT NULL REFERENCES works(id) ON DELETE CASCADE,
    child_id UUID REFERENCES work_children(id) ON DELETE CASCADE,
    source TEXT NOT NULL CHECK (source IN ('instance', 'local_file', 'public_domain', 'official_link')),
    locator TEXT NOT NULL,
    quality TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    -- Upsert key: re-scanning a library reconciles instead of duplicating.
    UNIQUE NULLS NOT DISTINCT (user_id, work_id, child_id, source, locator)
);

CREATE INDEX IF NOT EXISTS idx_content_references_resolve ON content_references(user_id, work_id, child_id);

CREATE OR REPLACE TRIGGER update_works_updated_at BEFORE UPDATE ON works FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE OR REPLACE TRIGGER update_work_children_updated_at BEFORE UPDATE ON work_children FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE OR REPLACE TRIGGER update_content_references_updated_at BEFORE UPDATE ON content_references FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
