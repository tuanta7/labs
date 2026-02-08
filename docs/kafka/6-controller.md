# Kafka Controller

In modern Kafka (versions 2.8+ using KRaft mode, and mandatory in 4.0+), the dependency on ZooKeeper has been eliminated. The controller role has been enhanced to a quorum of controllers that use the Raft consensus protocol to manage metadata internally within Kafka itself. 