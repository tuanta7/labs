# Explicit Locking

Locking is the tool, and isolation is the goal it helps achieve, with higher isolation levels typically requiring more complex locking

Most PostgreSQL commands automatically acquire locks of appropriate modes to ensure that referenced tables are not dropped or modified in incompatible ways while the command executes. There are majorly 2 types of locks:

- **Shared Locks** (Read Locks): Allow reading the row or the table that is being locked
- **Exclusive Locks**: Lock the row or table entirely and let the transaction update the row in isolation.

## 1. Table-Level Locks

## 2. Row-Level Locks

## Appendix: Optimistic & Pessimistic Locking

Optimistic locking is primarily managed within the application layer, though it relies on database support for concurrency control.
