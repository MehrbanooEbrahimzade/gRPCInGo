package main

import (
	"context"
	"log"
	"math/rand"
	"net"

	pb "github.com/MehrbanooEbrahimzade/gRPCInGo/users"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v", in.GetUserName())

	return &pb.User{ID: int32(rand.Intn(100)), UserName: in.GetUserName(), Email: in.GetEmail(), MobileNo: in.GetMobileNo(),
		BirthDate: in.GetBirthDate(), Password: in.GetPassword(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
