# Using Interceptor

## Run Server and Send Request

```sh
tailmoon:interceptor asa-taka$ go run main.go
```

to start server and send requests

```sh
grpc_cli call localhost:10000 Hello "name: 'asa-taka'"
```

from another terminal, then the interceptor message will be displayed.

```
2018/09/13 22:20:39 gRPC server starts on localhost:10000
2018/09/13 22:20:41 My interceptor called!
2018/09/13 22:20:41 /hello.Greeting/Hello: name:"asa-taka"  -> message:"Hello asa-taka, I am gentle-server"
```

## Note

### Generate Go Stub

```sh
protoc -I proto --go_out=plugins=grpc:proto hello.proto
```
