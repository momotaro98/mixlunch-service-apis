# db

## init

Initial script

## schema_changes

Diff changes for schema

## erd

Entity Relationship Diagram by using PlantUML

## mysql

MySQL configuration files

## Flow for updating DB schema and managing it

### Local development

1. Update `01_ddl.sql` and `02_dml.sql` scripts in `init` directory and check if it's OK in local docker-compose MySQL

2. Run

```
$ ./schema_change_flow.sh master > schema_changes/XXXX.sql
```

master means target branch  
then check if `XXXX.sql` file was generated in `schema_changes` directory

### Apply to STG/Prod DB

Apply the new script in `schema_changes` to (STG/PROD) MySQL by using
CUI like `$ mysql -hIPAddress -uUsername -p'pass' mixlunch < XXXX.sql`  
or  
GUI

### Update schema views

```
$ direnv allow
$ make docker-mid-run
$ cd db
$ docker run -v "$PWD/schema_view:/output" --net="host" schemaspy/schemaspy:6.1.0 -t mysql -host $DB_HOST:$DB_PORT -db $DB_DATABASE -u $DB_USER -p $DB_PASS -s $DB_DATABASE
```

Check by `$ cd schema_view` and `$ open .` and open with `index.html` with Web browser
