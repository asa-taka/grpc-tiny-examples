# Using Interceptor

Here is a simple example to implement interceptors of [grpc].

[grpc]: https://godoc.org/google.golang.org/grpc

First, define a function satisfies the interface of [grpc.UnaryServerInterceptor].

```go
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	res, err := handler(ctx, req)
	log.Println("My interceptor called!")
	log.Printf("%s: %v -> %v", info.FullMethod, req, res)
	return res, err
}
```

Then set it up to the `grpc.Server`.

```go
grpcServer := grpc.NewServer(
  grpc.UnaryInterceptor(loggingInterceptor),
)
```

Above function `loggingInterceptor` satisfy the interface of [grpc.UnaryServerInterceptor] and it can be converted to `ServerOption` by [grpc.UnaryInterceptor].

[grpc.UnaryInterceptor]: https://godoc.org/google.golang.org/grpc#UnaryInterceptor
[grpc.UnaryServerInterceptor]: https://godoc.org/google.golang.org/grpc#UnaryServerInterceptor


## Run Server and Send Request

```sh
go run main.go
```

to start server, and send requests

```sh
grpc_cli call localhost:10000 Hello "name: 'asa-taka'"
```

from another terminal, then the interceptor message will be displayed like below.

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
