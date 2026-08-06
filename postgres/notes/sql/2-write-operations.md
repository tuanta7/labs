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
