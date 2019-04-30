package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/siredmar/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("Greet request: ", req)
	firstname := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstname + " times " + strconv.Itoa(i)
		response := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(response)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

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
