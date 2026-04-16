/*
1. Descending Index

Use when
*/

DROP INDEX IF EXISTS idx_users_status_created_at;

-- No index
EXPLAIN SELECT * FROM users WHERE status='active' ORDER BY created_at DESC LIMIT 20;

--
CREATE INDEX idx_users_status_created_at ON users(status, created_at DESC);

-- With index
EXPLAIN SELECT * FROM users WHERE status='active' ORDER BY created_at DESC LIMIT 20;

DROP INDEX IF EXISTS idx_users_status_created_at;

/*
2. Partial Index

Use when
*/
