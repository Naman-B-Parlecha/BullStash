# Database Backup CLI Utility - BullStash

A command-line utility for backing up and restoring various types of databases. Built with Golang, this tool supports multiple DBMS, cloud storage, Discord notifications, and monitoring with Prometheus and Grafana.

## Features
- **Supported DBMS**: MySQL, PostgreSQL, MongoDB, SQLite.
- **Backup Types**: Full, incremental, and differential backups.
- **Storage Options**: Local and cloud (AWS S3, Google Cloud Storage, Azure Blob Storage).
- **Notifications**: Discord webhook integration.
- **Monitoring**: Prometheus and Grafana for real-time tracking.

## Commands and Flags
### 1. Backup Command
```bash
BullStash backup [flags]
```
**Flags:**
- `--dbtype`: Type of database (e.g., mysql, postgres, mongodb, sqlite).
- `--host`: Database host (default: localhost).
- `--port`: Database port (default: 3306 for MySQL, 5432 for PostgreSQL, etc.).
- `--user`: Database username.
- `--password`: Database password.
- `--dbname`: Database name.
- `--backup-type`: Type of backup (full, incremental, differential).
- `--output`: Output path for the backup file (default: `./backup.sql`).
- `--compress`: Enable compression (e.g., gzip) for the backup file.
- `--storage`: Storage option (local, s3, gcs, azure).
- `--cloud-bucket`: Cloud storage bucket name (required for cloud storage).
- `--cloud-region`: Cloud storage region (required for cloud storage).

**Example:**
```bash
BullStash backup --dbtype=postgres --host=localhost --port=5432 --user=myuser --password=mypassword --dbname=mydb --backup-type=full --output=./backup.sql --compress --storage=s3 --cloud-bucket=mybucket --cloud-region=us-east-1
```

### 2. Restore Command
```bash
BullStash restore [flags]
```
**Flags:**
- `--dbtype`: Type of database (e.g., mysql, postgres, mongodb, sqlite).
- `--host`: Database host (default: localhost).
- `--port`: Database port (default: 3306 for MySQL, 5432 for PostgreSQL, etc.).
- `--user`: Database username.
- `--password`: Database password.
- `--dbname`: Database name.
- `--input`: Path to the backup file (e.g., `./backup.sql`).
- `--storage`: Storage option (local, s3, gcs, azure).
- `--cloud-bucket`: Cloud storage bucket name (required for cloud storage).
- `--cloud-region`: Cloud storage region (required for cloud storage).
- `--selective`: Enable selective restore (e.g., specific tables or collections).

**Example:**
```bash
BullStash restore --dbtype=postgres --host=localhost --port=5432 --user=myuser --password=mypassword --dbname=mydb --input=./backup.sql --storage=s3 --cloud-bucket=mybucket --cloud-region=us-east-1
```

### 3. Schedule Command
```bash
BullStash schedule [flags]
```
**Flags:**
- `--dbtype`: Type of database (e.g., mysql, postgres, mongodb, sqlite).
- `--cron`: Cron expression for scheduling (e.g., `0 2 * * *` for daily at 2 AM).
- `--backup-type`: Type of backup (full, incremental, differential).
- `--storage`: Storage option (local, s3, gcs, azure).
- `--cloud-bucket`: Cloud storage bucket name (required for cloud storage).
- `--cloud-region`: Cloud storage region (required for cloud storage).

**Example:**
```bash
BullStash schedule --dbtype=mysql --cron="0 2 * * *" --backup-type=full --storage=s3 --cloud-bucket=mybucket --cloud-region=us-east-1
```

### 4. Test Command
```bash
BullStash test [flags]
```
**Flags:**
- `--dbtype`: Type of database (e.g., mysql, postgres, mongodb, sqlite).
- `--host`: Database host (default: localhost).
- `--port`: Database port (default: 3306 for MySQL, 5432 for PostgreSQL, etc.).
- `--user`: Database username.
- `--password`: Database password.
- `--dbname`: Database name.

**Example:**
```bash
BullStash test --dbtype=postgres --host=localhost --port=5432 --user=myuser --password=mypassword --dbname=mydb
```

### 5. Monitor Command
```bash
BullStash monitor [flags]
```
**Flags:**
- `--prometheus-port`: Port for exposing Prometheus metrics (default: 9090).
- `--grafana-dashboard`: Path to Grafana dashboard configuration file.

**Example:**
```bash
BullStash monitor --prometheus-port=9090 --grafana-dashboard=./grafana/dashboard.json
```

### 6. Notify Command
```bash
BullStash notify [flags]
```
**Flags:**
- `--discord-webhook`: Discord webhook URL for notifications.
- `--slack-webhook`: Slack webhook URL for notifications (optional).

**Example:**
```bash
BullStash notify --discord-webhook=https://discord.com/api/webhooks/your-webhook-url
```

---

## Milestones

### **Milestone 1: Core Functionality**
- [X] Implement basic CLI structure using Golang.
- [X] Add support for connecting PostgreSQL.
- [ ] Add support for connecting to MySQL
- [X] Implement full backup functionality for supported databases.
- [X] Add local storage option for backup files.
- [X] Implement basic error handling and logging.

### **Milestone 2: Advanced Backup Features**
- [ ] Add support for MongoDB and SQLite.
- [ ] Implement incremental and differential backups.
- [X] Add compression for backup files (e.g., using gzip).
- [ ] Implement connection testing for databases.
- [X] Implement Restore command

### **Milestone 3: Cloud Integration**
- [ ] Add support for cloud storage (AWS S3, Google Cloud Storage, Azure Blob Storage).
- [ ] Implement configuration file for storing cloud credentials and settings.
- [ ] Add functionality to schedule backups using cron-like syntax.

### **Milestone 4: Notifications and Monitoring**
- [ ] Integrate Discord webhooks for backup status notifications.
- [ ] Add Prometheus metrics for monitoring backup operations.
- [ ] Set up Grafana dashboards for visualizing backup performance and status.

### **Milestone 5: Restore Functionality**
- [ ] Implement restore functionality for all supported databases.
- [ ] Add selective restore options for tables or collections (if supported by the DBMS).
- [ ] Test restore operations with large databases for reliability.

### **Milestone 6: Optimization and Documentation**
- [ ] Optimize backup and restore operations for large databases.
- [ ] Add comprehensive documentation for CLI usage and configuration.
- [ ] Perform cross-platform testing (Windows, Linux, macOS).
- [ ] Write unit tests and integration tests for all features.

---

## Stretch Goals
- [ ] Add support for additional DBMS (e.g., Oracle, Cassandra).
- [ ] Implement encryption for backup files.
- [ ] Add support for Slack notifications alongside Discord.
- [ ] Create a Docker image for easy deployment.

---

#### Maybe add Encryption like Rishabh bhaiya thought us will be cool 

---