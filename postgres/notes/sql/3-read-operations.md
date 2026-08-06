## 3. Read Operations: SELECT

In SQL, read operations generally involve `SELECT` statements. The `SELECT` statement and its associated clauses are part of Data Manipulation Language (DML).

```sql
SELECT customer_name, city FROM customers;

-- select with alias
SELECT customer_id AS id FROM customers;
```

### 3.1. Clauses

#### ORDER BY

Sorts the result set by one or more columns.

```sql
SELECT * FROM table_name ORDER BY column1 ASC, column2 DESC;
```

#### LIMIT

Limits the number of rows returned by a query.

- `SELECT TOP` is implemented in MSSQL, and `LIMIT` is implemented in MySQL, MariaDB, and PostgreSQL, both are used to limit number of records returned

```sql
SELECT * FROM table_name LIMIT 10;
```

#### WHERE

Filters data based on certain conditions.

- Use `IS NULL` and `IS NOT NULL` instead of `=` or `!=` (or `<>`) when checking for `NULL` values
- Must always be used before the `GROUP BY` clause

```sql
SELECT * FROM employees WHERE middle_name IS NOT NULL;
SELECT name FROM customers WHERE referee_id IS NULL OR referee_id <> 2;
```

#### DISTINCT

Inside a table, a column often contains many duplicate values; and sometimes you only want to list the different values.

```sql
SELECT DISTINCT Country FROM Customers;

-- return the number of different countries
SELECT COUNT(DISTINCT Country) FROM Customers;
```

#### GROUP BY & HAVING

Aggregate and filter data based on specific groups.

- Allows you to group rows that have the same values in specified columns into summary row.
- Often combined with aggregate functions like COUNT, SUM, AVG, MIN, MAX, etc.
- You can group by multiple columns to create more granular groups.

```sql
SELECT category, SUM(amount) AS total_sales
FROM sales
GROUP BY category;
```

- Use `HAVING` to apply conditions on aggregate functions, because `WHERE` cannot filter by aggregate results.
- Example: Find categories where the average sale amount is over $500, but only include sales over $100 in the calculation.

```sql
SELECT category, AVG(amount) AS avg_sales
FROM sales
WHERE amount > 100
GROUP BY category
HAVING AVG(amount) > 500;
```

### 3.2. Aggregate Functions

An aggregate function is a function that performs a calculation on a set of values, and returns a single value.

- Often used with the `GROUP BY` clause of the `SELECT` statement
- Ignore NULL values, except for COUNT().
- The most commonly used SQL aggregate functions are COUNT(), SUM(), AVG(), MIN(), MAX().
