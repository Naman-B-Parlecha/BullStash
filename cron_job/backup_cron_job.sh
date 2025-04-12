#!/bin/bash
cd "/home/naman/Desktop/github/BullStash" || exit 1

echo "[$(date)] Starting BullStash backup..." >> "/home/naman/Desktop/github/BullStash/cron_logs/backup.log"

BullStash backup --dbtype postgres --backup-type full --output "backups" \
    >> "/home/naman/Desktop/github/BullStash/cron_logs/backup.log" 2>&1

echo "[$(date)] Backup completed." >> "/home/naman/Desktop/github/BullStash/cron_logs/backup.log"
