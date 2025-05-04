# QueryBuilder Package

The `QueryBuilder` package provides a flexible and intuitive way to construct SQL queries programmatically in Go. It is designed to simplify the process of building complex queries by offering a chainable API for defining query components such as fields, conditions, joins, and Common Table Expressions (CTEs).

## Features

- **Dynamic Query Building**: Construct SQL queries dynamically with support for SELECT, WHERE, JOIN, and CTEs.
- **Chainable API**: Easily chain methods to build queries step by step.
- **Parameterized Queries**: Automatically generates parameterized queries to prevent SQL injection.

## Usage

### Creating a QueryBuilder Instance
To start building a query, create a new `QueryBuilder` instance by specifying the table name:

```go
qb := querybuilder.NewQueryBuilder("users")
```

### Selecting Fields
Specify the fields to include in the SELECT statement:

```go
qb.Select("id", "name", "email")
```

### Adding Conditions
Add conditions to the WHERE clause using `Where` and `OrWhere` methods:

```go
qb.Where("age", ">", 18).OrWhere("status", "=", "active")
```

### Adding Joins
Define relationships between tables using the `Join` method:

```go
qb.Join("orders", "id", "user_id")
```

### Adding Common Table Expressions (CTEs)
Include CTEs in your query using the `AddCTE` method:

```go
cte := querybuilder.NewQueryBuilder("recent_orders").Select("id", "user_id").Where("created_at", ">", "2025-01-01")
qb.AddCTE("recent_orders_cte", cte)
```

### Building the Query
Generate the final SQL query and its arguments:

```go
query, args := qb.BuildQuery()
fmt.Println("Query:", query)
fmt.Println("Arguments:", args)
```

### Example
Here is a complete example of building a query:

```go
qb := querybuilder.NewQueryBuilder("users")
qb.Select("users.id", "users.name", "orders.total")
  .Join("orders", "id", "user_id")
  .Where("users.age", ">", 18)
  .OrWhere("users.status", "=", "active")

query, args := qb.BuildQuery()
fmt.Println("Query:", query)
fmt.Println("Arguments:", args)
```

## API Reference

### `NewQueryBuilder(tableName string) *QueryBuilder`
Creates a new `QueryBuilder` instance for the specified table.

### `Select(fields ...string) *QueryBuilder`
Specifies the fields to include in the SELECT statement.

### `Where(field, operator string, value interface{}) *QueryBuilder`
Adds a condition to the WHERE clause.

### `OrWhere(field, operator string, value interface{}) *QueryBuilder`
Adds an OR condition to the WHERE clause.

### `Join(table, foreignKey, primaryKey string) *QueryBuilder`
Adds a JOIN clause to the query.

### `AddCTE(name string, queryBuilder *QueryBuilder) *QueryBuilder`
Adds a Common Table Expression (CTE) to the query.

### `BuildQuery() (string, []interface{})`
Generates the final SQL query and its arguments.

## License
This package is open-source and available under the MIT License.
