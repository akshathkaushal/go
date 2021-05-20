package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	test "test/service"
	"time"
)

type clientUser struct {
	ID int64
	Name string
}

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err!=nil {
		panic(err)
	}

	defer conn.Close()
	client := test.NewDoOperationsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// adding users
	r1, err := client.Add(ctx, &test.UserData{ID: 1, Name: "Akshath"})
	if err != nil {
		panic(err)
	}
	log.Println(r1.GetID(), r1.GetName())

	r2, err := client.Add(ctx, &test.UserData{ID: 2, Name: "Kaushal"})
	if err != nil {
		panic(err)
	}
	log.Println(r2.GetID(), r2.GetName())

	// retrieving users

	r, err := client.Get(ctx, &test.GetUser{ID: 1})
	if err != nil {
		panic(err)
	}
	log.Println(r.GetID(), r.GetName())




}
