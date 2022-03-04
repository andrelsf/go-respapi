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

## References

- [Simple Bank](https://github.com/techschool/simplebank)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)