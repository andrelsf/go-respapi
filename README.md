# Golang REST API


## Golang Migrate

Install linux
```shell
curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
sudo echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/migrate.list
apt-get update
apt-get install -y migrate
```


## References

- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)