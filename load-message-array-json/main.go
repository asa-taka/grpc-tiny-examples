package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/asa-taka/grpc-tiny-examples/load-message-array-json/proto"
	"github.com/golang/protobuf/jsonpb"
)

const (
	port     = 10000
	dataFile = "data/todos.json"
)

// gRPC Service Implementation
// ---------------------------

type todoServer struct {
	todos []*pb.Todo
	mu    sync.Mutex
}

func newServer() *todoServer {
	s := &todoServer{}

	// Working pattern
	// s.loadTodos()

	// Practice: Try bad pattenrs and confirm results instead of `loadTodos()`
	s.loadTodos_byStandardJSONPackage()

	return s
}

// Load Array of `grpc.Message` from single JSON file.
// This technique is introduced on [golang/protobuf#675].
// [golang/protobuf#675]: https://github.com/golang/protobuf/issues/675#issuecomment-411131669
func (s *todoServer) loadTodos() {

	jsonBytes, _ := ioutil.ReadFile(dataFile)
	jsonString := string(jsonBytes)
	jsonDecoder := json.NewDecoder(strings.NewReader(jsonString))

	// read open bracket
	if _, err := jsonDecoder.Token(); err != nil {
		log.Fatal(err)
	}

	for jsonDecoder.More() {
		todo := pb.Todo{}
		if err := jsonpb.UnmarshalNext(jsonDecoder, &todo); err != nil {
			log.Fatal(err)
		}
		s.todos = append(s.todos, &todo)
	}
}

// Bad Pattern: Using Array of Message type directly as the Unmarshaling target
//
// It causes the error
//
// ```
// []*todo.Todo does not implement "github.com/gogo/protobuf/proto".Message
// (missing ProtoMessage method)
// ```
//
// Note: Commented out because of compile error
func (s *todoServer) loadTodos_byInstinctiveJSONPBUsage() {

	// reader, _ := os.Open(dataFile)

	// if err := jsonpb.Unmarshal(reader, &s.todos); err != nil {
	// 	log.Fatal(err)
	// }
}

// Bad Pattern: Using standard JSON library
//
// It causes the error
//
// ```
// json: cannot unmarshal string into Go struct field Todo.deadline of type timestamp.Timestamp
// ```
func (s *todoServer) loadTodos_byStandardJSONPackage() {

	jsonBytes, _ := ioutil.ReadFile(dataFile)

	if err := json.Unmarshal(jsonBytes, &s.todos); err != nil {
		log.Fatal(err)
	}
}

func (s *todoServer) GetTodos(ctx context.Context, in *pb.GetTodosRequest) (*pb.GetTodosResponse, error) {
	res := &pb.GetTodosResponse{
		Todos: s.todos,
	}
	return res, nil
}

// Main Program
// ------------

func main() {

	grpcServer := grpc.NewServer()

	// Register service implementations
	pb.RegisterReadOnlyTodoServer(grpcServer, newServer())

	reflection.Register(grpcServer) // for grpc_cli

	// ...then, followings are basic gRPC server setup...

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server starts on localhost:%d", port)
	grpcServer.Serve(lis)
}
