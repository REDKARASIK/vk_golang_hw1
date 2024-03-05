# README

<https://github.com/golang-standards/project-layout/blob/master/README_ru.md>

go mod init myapp

// go mod init github.com/rvasily/myapp

go mod download
go mod verify
go mod tidy

go build  -o ./bin/myapp ./cmd/myapp
go test -v -coverpkg=./... ./...

go mod vendor
go build -mod=vendor -o ./bin/myapp ./cmd/myapp
go test -v -mod=vendor -coverpkg=./... ./...
