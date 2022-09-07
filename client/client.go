package main

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"

	pb "github.com/MehrbanooEbrahimzade/gRPCInGo/users"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

type newUser struct {
	userName  string
	email     string
	mobileNo  string
	birthDate time.Time
	password  string
}

func main() {

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newUsers := []pb.User{
		{
			UserName:  "Ali",
			Email:     "ali.com",
			MobileNo:  "0915",
			BirthDate: "2009-11-17",
			Password:  "1234",
		}, {
			UserName:  "ahmad",
			Email:     "ahmad.com",
			MobileNo:  "0903",
			BirthDate: "2002-03-10",
			Password:  "9876",
		},
	}

	for _, u := range newUsers {
		r, err := c.CreateUser(ctx, &pb.CreateUserReq{User: &u})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf("NAME: %s", r.User.GetUserName())
	}
}
