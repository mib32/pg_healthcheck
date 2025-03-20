Simple healthcheck deamon for postgres.

`/health` endpoint calls `SELECT 1;` on the database.

Configuration options:
`--web.listenAddress`
`--dsn` - connection string

```
go build -o pg_healthcheck main.go
```