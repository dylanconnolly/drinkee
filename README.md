go/gin backend service

### Migrations

Postgres with sqlx + migrate

To create new migration
```
migrate create -ext sql -dir db/migrations -seq <migration_name>
```
To run migration
```
migrate -database "postgres://localhost:5432/drinkee?sslmode=disable" -path db/migrations up
```
To force DB to version
```
migrate -database "postgres://localhost:5432/drinkee?sslmode=disable" -path db/migrations force <version>
```
e.g. `migrate -database "postgres://localhost:5432/drinkee?sslmode=disable" -path db/migrations force 10`


### Testing
Integration tests use dockertest to spin up and tear down containers running instances of postgres DBs.

Test a single file
```
go test -race /path/to/file
```

for verbose output on tests:
```
go test /path/to/file -v
```
