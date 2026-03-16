# Concurrency Control in PostgreSQL

Reference: [PostgreSQL | Concurrency Control](https://www.postgresql.org/docs/current/mvcc.html)

Data consistency in PostgreSQL is maintained by using a multiversion model (Multiversion Concurrency Control, [MVCC](https://www.postgresql.org/docs/7.1/mvcc.html)). Each SQL statement sees a snapshot of data (a database version) as it was some time ago, regardless of the current state of the underlying data. This implementation provides transaction isolation for each database session.

> [!NOTE]
> A transaction is effectively treated as concurrent with another if the other transaction appears active in the snapshot.

## 1. Transaction Isolation (Implicit Locking)

A database transaction is a unit of work that is performed on a database and treated as a single atomic operation. It ensures that either all the changes within the transaction are committed to the database or none of them are. This guarantees data consistency and integrity.

The SQL standard defines four levels of transaction isolation. The most strict is Serializable, which guarantee that any concurrent execution of a set of transactions will produce the same effect as running them one at a time in some order.

The other three levels are defined in terms of phenomena, resulting from interaction between concurrent transactions, which must not occur at each level

- **Dirty Read**: A transaction reads data written by a concurrent uncommitted transaction.
- **Non-repeatable Read**: A transaction re-reads data it has previously read and finds that data **VALUE** has been modified by another recently-committed transaction.
- **Phantom Read**: A transaction re-executes a query returning a set of rows that satisfy a search condition and finds that the **SET** of rows satisfying the condition has changed due to another recently-committed transaction.
- **Serialization Anomaly**: The result of successfully committing a group of transactions is inconsistent with **ALL** possible orderings of running those transactions one at a time.

| Isolation Level  | Dirty Read             | Nonrepeatable Read | Phantom Read           | Serialization Anomaly | _Count_ |
| ---------------- | ---------------------- | ------------------ | ---------------------- | --------------------- | ------- |
| Read uncommitted | Allowed, but not in PG | Possible           | Possible               | Possible              | 3       |
| Read committed   | Not possible           | Possible           | Possible               | Possible              | 3       |
| Repeatable read  | Not possible           | Not possible       | Allowed, but not in PG | Possible              | 1       |
| Serializable     | Not possible           | Not possible       | Not possible           | Not possible          | 0       |

> [!NOTE]
> Each transaction is sequential and isolated. Within the same transaction, PostgreSQL does not ignore prior changes from earlier statements. The second update/read of a row always sees the first update's result.

To set the transaction isolation level of a transaction (default to READ COMMITTED), use the command [SET TRANSACTION](https://www.postgresql.org/docs/current/sql-set-transaction.html). Note that internally only three distinct isolation levels are implemented (Read Uncommitted mode behaves like Read Committed)

### 1.1. Read Committed Isolation Level

Read Committed mode starts each command with a new snapshot that includes all transactions committed up to that instant

When a transaction uses this isolation level, a SELECT query (without a FOR UPDATE/SHARE clause) sees only data committed before the query began; it never sees either uncommitted data or changes committed by concurrent transactions during the query's execution.

```sql
-- Transaction #1
BEGIN;
UPDATE accounts SET balance = balance + 100 WHERE acctnum = 123;
-- Balance is now 200 in this session, but not yet committed.
SELECT balance FROM accounts WHERE acctnum = 123; -- Returns 200

-- Transaction #2 (concurrently)
SELECT balance FROM accounts WHERE acctnum = 123; -- Returns the committed balance (e.g., 100)

-- Transaction #1
COMMIT;

-- Transaction #2 (concurrently)
SELECT balance FROM accounts WHERE acctnum = 123; -- Now returns 200
```

### 1.2. Repeatable Read Isolation Level

This level is different from Read Committed in that a query in a repeatable read transaction sees a **STABLE** snapshot at the beginning of the transaction

> [!NOTE]
> A non-transaction-control statement is any SQL statement that isn't a transaction control statement (like COMMIT, ROLLBACK, or SAVEPOINT) and manipulates data or defines data structures. For example: SELECT, INSERT, UPDATE, DELETE. CREATE, ALTER, DROP, etc.

```sql
-- Transaction #1
BEGIN TRANSACTION ISOLATION LEVEL REPEATABLE READ;
SELECT COUNT(*) FROM users WHERE status = 'active'; -- Returns, for example, 5

-- Transaction #2 (concurrently)
BEGIN;
INSERT INTO users (name, status) VALUES ('New User', 'active');
COMMIT; -- new change is commited

-- Transaction #1
SELECT COUNT(*) FROM users WHERE status = 'active'; -- Still returns 5, even though a new user was added and committed.
COMMIT;
```

Applications using this level must be prepared to retry transactions due to serialization failures. If the first updater commits (and actually updated or deleted the row, not just locked it) then the repeatable read transaction will be rolled back with the message

```shell
ERROR:  could not serialize access due to concurrent update
```

because a repeatable read transaction cannot modify or lock rows changed by other transactions after the repeatable read transaction began.

### 1.3. Serializable Isolation Level

This level emulates serial transaction execution for all committed transactions; as if transactions had been executed one after another, serially, rather than concurrently. It works exactly the same as Repeatable Read except that it also monitors for conditions which could make execution of a concurrent set of serializable transactions behave in a manner inconsistent with all possible serial executions of those transactions.

```sql
-- Transaction #1: Alice checks if she can leave duty
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
SELECT COUNT(*) FROM on_call WHERE on_call = true;  -- result = 2
UPDATE on_call SET on_call = false WHERE doctor = 'Alice';
COMMIT;

-- Transaction #2 (concurrently): Bob performs the same check
BEGIN TRANSACTION ISOLATION LEVEL SERIALIZABLE;
SELECT COUNT(*) FROM on_call WHERE on_call = true;  -- result = 2
UPDATE on_call SET on_call = false WHERE doctor = 'Bob';
COMMIT;

-- The constraint is violated because no doctor remains on call.
ERROR:  could not serialize access due to read/write dependencies among transactions
```

If these transactions were not serializable, and one transaction inserted/updated a row that would affect the other's calculation, a non-serializable outcome could occur.

## 2. Explicit Locking

Locking is the tool, and isolation is the goal it helps achieve, with higher isolation levels typically requiring more complex locking

Most PostgreSQL commands automatically acquire locks of appropriate modes to ensure that referenced tables are not dropped or modified in incompatible ways while the command executes.

There are majorly 2 types of locks:

- **Shared Locks** (or read lock): Allow reading the row or the table that is being locked
- **Exclusive Locks**: Lock the row or table entirely and let the transaction update the row in isolation.

### 2.1. Table-Level Locks

### 2.2. Row-Level Locks

## 3. Optimistic & Pessimistic Locking

Optimistic locking is primarily managed within the application layer, though it relies on database support for concurrency control.
