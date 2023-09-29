CREATE OR REPLACE FUNCTION all_jobs() RETURNS TABLE (
        id BIGINT,
        job_id VARCHAR(36),
        message VARCHAR(500),
        created_at TIMESTAMP,
        output JSON
    ) AS $$
SELECT id,
    job_id,
    message,
    created_at,
    output
FROM jobs
ORDER BY created_at DESC;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION all_workers() RETURNS TABLE (
        id VARCHAR(36),
        name VARCHAR(255),
        description VARCHAR(255)
    ) AS $$
SELECT id,
    name,
    description
FROM workers 
$$ LANGUAGE SQL;
