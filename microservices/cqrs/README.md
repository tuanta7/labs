# Command Query Responsibility Segregation (CQRS)

Reference: [Azure Architecture Center](https://learn.microsoft.com/en-us/azure/architecture/patterns/cqrs)

An architectural pattern that separates the concerns of reading and writing data. It divides an application into two distinct parts:

- **Command Side**: Responsible for managing create, update, and delete requests.
- **Query Side:** Responsible for handling read requests.
