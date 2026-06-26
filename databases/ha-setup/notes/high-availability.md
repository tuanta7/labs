# HA Systems

- **Active-Passive Load Balancing**: For standard web servers or applications, a main server handles traffic while a secondary server waits on standby. When the main server fails, traffic simply routes to the backup.
- **Stateless Microservices** (Redundancy): Running multiple copies at once.
- **Distributed Databases** (Replication): Use a consensus algorithm to ensure that multiple nodes agree on transaction history, which in turn maintains a continuously available.
