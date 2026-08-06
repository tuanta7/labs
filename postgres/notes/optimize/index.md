# Optimizing SQL Queries

## 1. Indexing

Indexes are auxiliary data structures that allow the database to locate rows more quickly, avoiding full table scans when possible. They function like a book’s index, pointing to relevant pages (rows) instead of requiring a complete read.

```sql
-- Create an index on the id column
CREATE INDEX user_id_index ON users (id);

-- Remove an index
DROP INDEX user_id_index;
```

Indexes are beneficial for queries with WHERE clauses, JOIN conditions, and ORDER BY operations. However, their use depends on the query planner's cost analysis. An index may be omitted if:

- The dataset is very small (a sequential scan may be faster).
- The query requests a large portion of the table (scanning is cheaper than many index lookups).
- The index does not align with the filtering or sorting conditions.

```sql
EXPLAIN ANALYZE SELECT full_name FROM account WHERE mobile = '097****731';
```

```sh
Seq Scan on account  (cost=0.00..2507.64 rows=1 width=23) (actual time=0.283..1.994 rows=1 loops=1)
  Filter: ((mobile)::text = '097****731'::text)
  Rows Removed by Filter: 956
Planning Time: 0.154 ms
Execution Time: 2.025 ms
```

```sql
EXPLAIN ANALYZE SELECT full_name FROM account WHERE account_no = '00123456'; -- with index
```

```sh
Index Scan using account_account_no_idx on account  (cost=0.40..2.62 rows=1 width=23) (actual time=0.060..0.061 rows=1 loops=1)
  Index Cond: ((account_no)::text = '00123456'::text)
Planning Time: 0.201 ms
Execution Time: 0.090 ms
```

### Trade-offs of Indexing

- **Storage Overhead**: Indexes consume additional disk space.
- **Write Cost**: INSERT, UPDATE, and DELETE operations become slower since indexes must be updated.
- **Over-indexing**: Too many indexes can degrade overall performance.

Beyond indexing and execution plan analysis, several techniques can further improve SQL performance:

- **Query Rewriting**: Restructure queries to make them more optimizer-friendly (e.g., replacing subqueries with joins).
- **Partitioning**: Splitting large tables into smaller, more manageable chunks.
- **Materialized Views**: Storing precomputed query results to avoid expensive recalculations.
- **Caching**: Using in-memory caches for frequently accessed results.
- **Vacuuming** (PostgreSQL): Reclaiming storage and maintaining healthy table structures.
