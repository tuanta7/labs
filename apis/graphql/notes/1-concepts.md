# GraphQL

Reference: [graphql.org](https://graphql.org/learn/introduction/)

GraphQL is a query language for your API, and a server-side runtime for executing queries using a type system.

## 1. Schemas & Types

The GraphQL type system describes what data can be queried from the API. The collection of those capabilities is referred to as the service’s schema and clients can use that schema to send queries to the API that return predictable results.

The most basic components of a GraphQL schema are Object types, which just represent a kind of object that can be fetched from a service, and what fields it has.

```gql
type Character {
  name: String!
  appearsIn: [Episode!]!
}
```

Every field on a GraphQL Object type can have zero or more arguments

- All arguments are named.

```gql
type Starship {
  id: ID!
  name: String!
  length(unit: LengthUnit = METER): Float
}
```

The GraphQL query language is basically about selecting fields on objects.

```txt

```

### Validation

In practice, when a GraphQL operation reaches the server, the document is first parsed and then validated using the type system. This allows servers and clients to effectively inform developers when an invalid query has been created, without relying on runtime checks.

## 2. Queries, Mutations & Subscriptions

Every GraphQL schema must support query operations.
