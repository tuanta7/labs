# Minimal Monitoring Setup

| Component               | Role                                                                                                  |
| ----------------------- | ----------------------------------------------------------------------------------------------------- |
| OpenTelemetry SDKs      | Instrument services with automatic or manual traces, logs and metrics.                                |
| OpenTelemetry Collector | Collects metrics, traces, logs; forwards to SigNoz Collector                                          |
| SigNoz (All-in-one)     | Observability platform native to OpenTelemetry with logs, traces and metrics in a single application. |
| ClickHouse              | A lightning-fast, column-oriented database engineered specifically for analytical workloads           |

## HTTPS

- Certificate Authority (CA): A trusted entity that issues digital certificates. Examples include Let's Encrypt, DigiCert, and Comodo.
- Digital Certificate: A file that contains a public key and a digital signature.
- Private Key: A secret key that is used to sign a digital certificate.
- TLS: A secure protocol that encrypts data in transit between a client and a server.
- SSL (deprecated): An older protocol that has been replaced by TLS.

## APT (Advanced Package Tool)

`apt` and `apt-get` are components of the APT (Advanced Package Tool) ecosystem used on Debian and Ubuntu distributions.

### Repository Index Retrieval

When apt update or apt-get update is executed, the following sequence occurs

- Entries in `/etc/apt/sources.list` and the files under `/etc/apt/sources.list.d/` are parsed.
