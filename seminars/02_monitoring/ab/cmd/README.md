```shell
go build -o ab ./...
./ab -k -c 8 -n 100000 http://localhost:8080/hello/world
```
