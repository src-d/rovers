CREATE TABLE IF NOT EXISTS repositories (
	id uuid PRIMARY KEY,
	created_at timestamptz,
        updated_at timestamptz,
        endpoints text[],
        status varchar(20),
        fetched_at timestamptz,
        fetch_error_at timestamptz,
        last_commit_at timestamptz,
        is_fork boolean,
        _references jsonb
);

CREATE INDEX IF NOT EXISTS idx_repositories_endpoints on "repositories" USING GIN ("endpoints");
