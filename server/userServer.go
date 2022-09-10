package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	pb "github.com/MehrbanooEbrahimzade/gRPCInGo/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	port = "localhost:50051"
)

var (
	client *mongo.Client
	Userdb *mongo.Collection
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(s, &UserServiceServer{})
	log.Printf("server listening at %v", lis.Addr())
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	configMongodb(ctx)

	fmt.Printf("Server succesfully started on port :%v", port)
	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("\nStopping the server...")
	s.Stop()
	lis.Close()
	fmt.Println("Closing MongoDB connection")
	client.Disconnect(ctx)
	fmt.Println("Done.")
}

func configMongodb(ctx context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to Mongodb")
	}

	Userdb = client.Database("mydb").Collection("user")

	mod := mongo.IndexModel{Keys: bson.M{"mobileNo": 1}, Options: options.Index().SetUnique(true)}
	Userdb.Indexes().CreateOne(ctx, mod)
}

type UserItem struct {
	ID        int32  `bson:"_id,omitempty"`
	UserName  string `bson:"userName"`
	Email     string `bson:"email"`
	MobileNo  string `bson:"mobileNo"`
	BirthDate string `bson:"birthDate"`
	Password  string `bson:"password"`
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRes, error) {
	fmt.Println("into create")
	user := req.GetUser()
	// convert into BSON
	data := UserItem{
		UserName:  user.UserName,
		Email:     user.Email,
		MobileNo:  user.MobileNo,
		BirthDate: user.BirthDate,
		Password:  user.Password,
	}
	fmt.Println("into create")

	result, err := Userdb.InsertOne(ctx, data)
	if err != nil {
		// return internal gRPC error to be handled later
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid := result.InsertedID.(int32)
	fmt.Println("into create")
	user.ID = oid
	return &pb.CreateUserRes{User: user}, nil
}
