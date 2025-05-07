## postgres 

```bash
BullStash backup --dbtype=postgres --host=localhost --port=5432 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_postgres --backup-type=full --output=backups --storage=cloud
```

```bash
BullStash restore --dbtype=postgres --host=localhost --port=5432 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_postgres --input=backups/bullstash_postgres_backup_20250508_031356.sql
```

```bash
BullStash test --dbtype=postgres --host=localhost --port=5432 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_postgres
```

```bash
BullStash schedule --dbtype=postgres --backuptype=full --output=backups --cron="* * * * *"
```

```bash
BullStash postgres --dbtype=postgres --host=localhost --port=5432 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_postgres
```


## mysql 

```bash
BullStash backup --dbtype=mysql --host=localhost --port=3306 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_mysql --backup-type=full --output=backups 
```

```bash
BullStash restore --dbtype=mysql --host=localhost --port=3306 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_mysql --input=backups/bullstash_mysql_backup_20250508_032247.sql
```

```bash
BullStash test --dbtype=mysql --host=localhost --port=3306 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_mysql
```

```bash
BullStash schedule --dbtype=mysql --backuptype=full --output=backups --cron="* * * * *"
```

```bash
BullStash mysql --dbtype=mysql --host=localhost --port=3306 --user=bullstash_user_1 --password=SpikyOP@123 --dbname=bullstash_mysql
```

## mongodb 

```bash
BullStash backup --dbtype=mongo --dbname=bullstash_mongo --backup-type=full --output=backups --compress=false --mongo-uri=mongodb+srv://bullstash_user_1:Si7y1zo1MjmNXVaI@bullstash-clluster-1.4ywveyu.mongodb.net/bullstash_mongo
```

```bash
BullStash restore --dbtype=mongo --mongo-uri=mongodb+srv://bullstash_user_1:Si7y1zo1MjmNXVaI@bullstash-clluster-1.4ywveyu.mongodb.net/bullstash_mongo --input=backups/backup_bullstash_mongo_20250508_032624 --drop=true --isCompressed=false
```

```bash
BullStash test --dbtype=mongo --mongo_uri=mongodb+srv://bullstash_user_1:Si7y1zo1MjmNXVaI@bullstash-clluster-1.4ywveyu.mongodb.net/bullstash_mongo 
```

```bash
BullStash schedule --dbtype=mongo --backuptype=full --output=backups --cron="* * * * *"
```

```bash
BullStash backup --dbtype=mongo --dbname=bullstash_mongo --mongo-uri=mongodb+srv://bullstash_user_1:Si7y1zo1MjmNXVaI@bullstash-clluster-1.4ywveyu.mongodb.net/bullstash_mongo
```

```bash
BullStash notify --discord=https://discord.com/api/webhooks/1361420028326711346/ikw4k92nKiXaDG6QFhToe858dbjwJXWTMk8u3uW5bvUnf_tdvUItXfQ4ZCetOPp9iz1_
```