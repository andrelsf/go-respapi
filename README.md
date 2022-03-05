# Golang REST API


## Golang Migrate

Install linux
```shell
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
sudo echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate
```


```shell
migrate -version
migrate -help
mkdir -p db/migration

migrate create -ext sql -dir db/migration -seq init_schema

migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
```

## SQLC

Install
```shell
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
sudo snap install sqlc

sqlc version
sqlc help

sqlc init
make sqlc
```

```shell
go mod init github.com/andrelsf/go-restapi
go mod tidy
```

```shell
go get github.com/lib/pq
go get github.com/stretchr/testify
```

Tests with HTTPie
```shell
http --json POST :8080/accounts < payloads/postAccounts.json
```

## References

- [Simple Bank](https://github.com/techschool/simplebank)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
- [SQLC](https://sqlc.dev/)
- [Github SQLC](https://github.com/kyleconroy/sqlc)
- [Go lib PQ](https://github.com/lib/pq)
- [Testify](https://github.com/stretchr/testify)
- [Golang GIN](https://github.com/gin-gonic/gin)