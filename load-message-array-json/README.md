# Load JSON Containing Array of Messages

Sometimes we want to load mock data from JSON files to serve
from such as first implementation of project's gRPC server.

In this example, **JSON containing Array of Message** is like below.

```json
[
  {
    "id": 1,
    "title": "First ToDo",
    "deadline": "2018-09-01T12:34:56.789Z"
  },
  {
    "id": 2,
    "title": "Second ToDo",
    "deadline": "2018-09-01T12:34:56.789Z"
  }
]
```

Thanks to [golang/protobuf#675], here succees to show working pattern.

[golang/protobuf#675]: https://github.com/golang/protobuf/issues/675#issuecomment-411131669

```sh
go run main.go
```

to start server, and send requests

```
$ grpc_cli call localhost:10000 GetTodos ""
connecting to localhost:10000
todos {
  id: 1
  title: "First ToDo"
  deadline {
    seconds: 1535805296
    nanos: 789000000
  }
}
todos {
  id: 2
  title: "Second ToDo"
  deadline {
    seconds: 1535805296
    nanos: 789000000
  }
}
```

## Note

### Generate Go Stub

```sh
protoc -I proto --go_out=plugins=grpc:proto todo.proto
```
