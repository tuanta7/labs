# Query Scenarios

## 1. Composite Indexes

Reference: [Composite vs Separate Indexes](https://www.cybertec-postgresql.com/en/combined-indexes-vs-separate-indexes-in-postgresql/)

A composite index works by creating a single, sorted data structure (usually a B-tree) based on the combined values of multiple columns, ordered first by the leftmost column, then the second, and so on

```sql
CREATE INDEX idx_data ON t_data(a, b, c);
```

### 1.1. Covering Indexes

Covering indexes are indexes that contain all the columns needed to satisfy a particular query without accessing the actual table data.

```sql
-- use INCLUDE to add extra columns (Postgres)
CREATE INDEX tab_x_y ON tab(x) INCLUDE (y);

-- y can be obtained from the index without visiting the heap
SELECT y FROM tab WHERE x = 'key';
```

### 1.2. Descending Indexes

Use the DESC keyword in the CREATE INDEX statement for the specific column to create a database index in descending order. This kind of index is effective for acquiring the top N latest rows that satisfy particular conditions (log/order/event).

```sql
users (
  id BIGINT PK,
  email VARCHAR,
  status VARCHAR,-- active, inactive
  created_at TIMESTAMP
)

CREATE INDEX idx_users_status_created_at ON users(status, created_at DESC);

-- get 20 newest active user
SELECT * FROM users
WHERE status='active'
ORDER BY created_at
DESC LIMIT 20;
```

### 1.3. Partial Indexes

Reference: [Partial Indexes](https://www.postgresql.org/docs/current/indexes-partial.html)

A partial index is an index built over a subset of a table; the subset is defined by a conditional expression. One major reason for using a partial index is to avoid indexing common values since a query searching for a common value will not use the index anyway.

```sql
tasks (
  id BIGINT PK,
  is_done BOOLEAN,
  created_at TIMESTAMP
)

CREATE INDEX idx_tasks_open_created ON tasks(created_at DESC) WHERE is_done = false;

-- get 50 newest open tasks
SELECT * FROM tasks
WHERE is_done = false
ORDER BY created_at DESC
LIMIT 50;
```

### 1.4. Joining with Filters

Applying a filter before joining usually makes joins faster by reducing the amount of data that needs to be shuffled, sorted, and compared.

```sql
users (
  id BIGINT PK,
  country VARCHAR
)

orders (
  id BIGINT PK,
  user_id BIGINT,
  status VARCHAR,
  created_at TIMESTAMP
)

--
SELECT o.*, u.* FROM orders AS o
JOIN users AS u
ON u.id = o.user_id
WHERE u.country = 'VN'
AND o.created_at >= NOW() - INTERVAL '7 days';
```

In absence of indexes, both users and orders are read via sequential scan.

- **users**: every row must be checked for country .
- **orders**: every row must be checked for created_at.

After filtering users, the join is executed without indexed lookup. It is often performed using Hash Join (most common, using extra memory for a hash map) or Nested Loop Join (worse case, for each matching user, the entire orders table is scanned).

```sql
-- fast retrieval of user IDs belonging to a country
CREATE INDEX idx_users_country_id ON users(country, id);

-- efficient lookup of recent orders per user
CREATE INDEX idx_orders_user_created ON orders(user_id, created_at);
```

Indexing makes filtering faster, then filtering make joining faster

### 1.5. Count

Index improves COUNT when a filter significantly reduces row set or an index-only scan is possible.

```sql

```

### 1.6. IN Operator

Use EXISTS and subqueries

### 1.7. OR Operator

```sql
SELECT * FROM orders
WHERE status = 'pending' OR created_at >= now() - interval '1 day';
```

Use UNION so planner can use indexes

```sql
 -- this column has low cardinality, a partial index may be better
CREATE INDEX idx_orders_status ON orders(status);

CREATE INDEX idx_orders_created ON orders(created_at);

SELECT * FROM orders WHERE status = 'pending'
UNION
SELECT * FROM orders WHERE created_at >= now() - interval '1 day';
```

### 1.8. Keyset Pagination (Cursor-based Pagination)

## 2. String

### LIKE
