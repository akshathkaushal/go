package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	test "test/service"
)

type server struct {
	test.UnimplementedDoOperationsServer
}

type user struct {
	ID   int64
	Name string
}

var userDataMap = map[int64]user{}

func (s *server) Add(ctx context.Context, inputData *test.UserData) (*test.UserData, error) {
	id, name := inputData.GetID(), inputData.GetName()

	var newUser user
	newUser.ID = id
	newUser.Name = name

	userDataMap[id] = newUser

	return &test.UserData{ID: newUser.ID, Name: newUser.Name}, nil
}

func (s *server) Get(ctx context.Context, inputID *test.GetUser) (*test.UserData, error) {
	id := inputID.GetID()

	if _, present := userDataMap[id]; present {
		return &test.UserData{ID: id, Name: userDataMap[id].Name}, nil
	}
	return nil, nil
}


func main() {
	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	test.RegisterDoOperationsServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(e)
	}
}