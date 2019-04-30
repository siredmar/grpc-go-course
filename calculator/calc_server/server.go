package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/siredmar/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calcpb.CalcRequest) (*calcpb.CalcResponse, error) {
	result := &calcpb.CalcResponse{
		Sum: req.First + req.Second,
	}
	return result, nil
}

func main() {
	fmt.Println("Running calc server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Could not create listener: %v", err)
	}
	s := grpc.NewServer()

	calcpb.RegisterCalcServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed serving: %v", err)
	}
}
