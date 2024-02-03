# toll-calculator

## Run Kafka Locally
```
    docker-compose up -d
```

## Installing protobuf compiler and setup

For linux

Note that you need to set the  /go/bin directory in your path 

```
    PATH="${PATH}:${HOME}/go/bin"
```

install protobuf and package dependencies

```
    sudo apt install -y protobuf-compiler

    go get google.golang.org/protobuf

    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

    go get google.golang.org/grpc

```