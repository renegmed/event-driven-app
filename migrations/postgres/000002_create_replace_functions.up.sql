CREATE OR REPLACE FUNCTION append_job_status(
        job_id VARCHAR(36),
        message VARCHAR(500),
        created_at TIMESTAMP,
        output JSON
    ) RETURNS VARCHAR(36) AS $$
INSERT INTO jobs(job_id, message, created_at, output)
VALUES (job_id, message, created_at, output)
RETURNING id;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION create_worker(
        name VARCHAR(255),
        description VARCHAR(255)
    ) RETURNS VARCHAR(36) AS $$
INSERT INTO workers(id, name, description)
VALUES (gen_random_uuid(), name, description)
RETURNING id;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION delete_worker(worker_id VARCHAR(36)) RETURNS VARCHAR(36) AS $$
DELETE FROM workers
WHERE id = worker_id
RETURNING id;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION get_worker(worker_id VARCHAR(36)) RETURNS TABLE (
        id VARCHAR(36),
        name VARCHAR(255),
        description VARCHAR(255)
    ) AS $$
SELECT id,
    name,
    description
FROM workers
WHERE id = worker_id
LIMIT 1;
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION job_statuses(j_id VARCHAR(36)) RETURNS TABLE (
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
WHERE job_id = j_id
ORDER BY created_at DESC;
$$ LANGUAGE SQL;
