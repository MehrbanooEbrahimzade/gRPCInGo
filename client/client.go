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
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	newUsers := []newUser{
		{
			userName:  "Ali",
			email:     "ali.com",
			mobileNo:  "0915",
			birthDate: time.Date(2009, 11, 17, 0, 0, 0, 0, time.UTC),
			password:  "1234",
		}, {
			userName:  "ahmad",
			email:     "ahmad.com",
			mobileNo:  "0903",
			birthDate: time.Date(2002, 03, 10, 0, 0, 0, 0, time.UTC),
			password:  "9876",
		},
	}

	for _, u := range newUsers {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{UserName: u.userName, Email: u.email, MobileNo: u.mobileNo,
			BirthDate: u.birthDate.Format(), Password: u.password})

		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf(`User Details:
		ID: %d
		NAME: %s
		EMAIL: %s
		MobileNo: %s
		BirthDate : %v`, r.GetID(), r.GetUserName(), r.GetEmail(), r.GetMobileNo(), r.GetBirthDate())
	}
}
