# SQL Basic

SQL (Structured Query Language) is the standard language for interacting with relational database management systems (RDBMS)

- There are formal standards for SQL (like ANSI SQL and ISO SQL) that define the core syntax and behavior.
- The SQL used in MariaDB, MySQL, or PostgreSQL are dialects of the same core language. Beyond standard SQL, these databases often incorporate their own procedural languages to support complex logic.

## 1. Data Types

### 1.1. Numeric

### 1.2. String

Reference: [MariaDB](https://mariadb.com/kb/en/string-data-types/)

- `CHAR`: A fixed-length string that is always right-padded with spaces to the specified length when stored
- `VARCHAR` (CHARACTER VARYING): A variable-length string that used for non-unicode (ASCII) characters only
- `NVARCHAR` (NATIONAL VARCHAR): A variable-length string that used for both unicode and non-unicode characters
- `TEXT`

### 1.3. Date and Time

### 1.4. Other Types

## 2. Create Operations: CREATE & INSERT

In SQL, `CREATE` operations are part of the Data Definition Language (DDL) commands, which are used to create or define database structures, such as tables, databases, indexes, and other objects.

```sql
CREATE TABLE hq_sales.invoices (
   invoice_id BIGINT UNSIGNED NOT NULL,
   branch_id INT NOT NULL,
   customer_id INT,
   invoice_date DATETIME(6),
   invoice_total DECIMAL(13, 2),
   payment_method ENUM('NONE', 'CASH', 'WIRE_TRANSFER', 'CREDIT_CARD', 'GIFT_CARD'),
   PRIMARY KEY (invoice_id)
);

CREATE INDEX idx_invoices_branch_customer ON hq_sales.invoices (branch_id, customer_id) IGNORED;
```

### 2.1. Field Constraints

#### UNIQUE

Requires values in column or columns only occur once in the table.

```sql

```

#### NOT NULL

Ensure that a column's value is not set to NULL

```sql

```

#### CHECK

Before a row is inserted or updated, all constraints are evaluated in the order they are defined.

- If any constraint expression returns false, then the row will not be inserted or updated.

```sql

```

#### PRIMARY KEY

Sets the column for referencing rows. Values must be `UNIQUE` and `NOT NULL`.

```sql
CREATE TABLE persons (
    id INT NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    age INT,
    PRIMARY KEY (id)
)
```

#### FOREIGN KEY

Sets the column to reference the primary key on another table.

```sql
CREATE TABLE orders (
    order_id INT NOT NULL,
    order_number INT NOT NULL,
    person_id INT,
    PRIMARY KEY (order_id),
    FOREIGN KEY (person_id) REFERENCES persons(id)
);
```

### 2.2. Insert

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

## 4. Modification Operations: ALTER, DROP & UPDATE, DELETE

## 5. Foreign Key: [MariaDB](https://mariadb.com/kb/en/foreign-keys/)

The `FOREIGN KEY` constraint is a key used to link two tables together. A `FOREIGN KEY` is a field (or collection of fields) in one table that refers to the `PRIMARY KEY` in another table.

```sql
CREATE TABLE Orders (
    OrderID int NOT NULL,
    OrderNumber int NOT NULL,
    PersonID int,
    PRIMARY KEY (OrderID),
    FOREIGN KEY (PersonID) REFERENCES Persons(PersonID)
);
```

### Key Constraints

> [!IMPORTANT]
> In MariaDB, the default behavior for both ON DELETE and ON UPDATE is `RESTRICT`.

#### OnDelete Constraints

This constraint specifies what happens to records in the child table when a record in the parent table is deleted.

| Options     | Description                                                                    |
| ----------- | ------------------------------------------------------------------------------ |
| CASCADE     | All related records in the child table are also deleted                        |
| SET NULL    | The change is allowed, and the child row's foreign key columns are set to NULL |
| SET DEFAULT | The foreign key in the child table is set to its default value.                |
| RESTRICT    | The change on the parent table is prevented                                    |
| NO ACTION   | Synonym for RESTRICT                                                           |

#### OnUpdate Constraints

| Options     | Description                                                                    |
| ----------- | ------------------------------------------------------------------------------ |
| CASCADE     | All related records in the child table are also updated (foreign key updated)  |
| SET NULL    | The change is allowed, and the child row's foreign key columns are set to NULL |
| SET DEFAULT | The foreign key in the child table is set to its default value.                |
| RESTRICT    | The change on the parent table is prevented                                    |
| NO ACTION   | Synonym for RESTRICT                                                           |

> [!TIP]
> In most database systems, ON DELETE CASCADE is common for ON DELETE, especially in one-to-many relationships where child records depend on the parent. ON UPDATE NO ACTION (or the default behavior) is common for ON UPDATE, as updating primary keys is less frequent and often discouraged.
