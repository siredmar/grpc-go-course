package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/siredmar/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetingRequest) (*greetpb.GreetingResponse, error) {
	fmt.Println("Greet request: ", req)
	firstname := req.GetGreeting().GetFirstName()
	result := "Hello " + firstname
	response := &greetpb.GreetingResponse{
		Result: result,
	}
	return response, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}