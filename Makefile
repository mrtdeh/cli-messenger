

.PHONY: all test clean

all: test build

test:
	go test ./...

build: clean
	go build -v -o bin/server  ./main.go
	go build -v -o bin/client  ./client/client.go

protoc:
	protoc --go-grpc_out=require_unimplemented_servers=false:./proto/ ./proto/*.proto --go_out=./proto




clean:
	rm -f ./bin/*