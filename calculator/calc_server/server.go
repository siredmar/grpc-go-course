package main

import (
	"context"
	"fmt"
	"io"
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

func (*server) Prime(req *calcpb.PrimeRequest, stream calcpb.Calc_PrimeServer) error {
	fmt.Printf("got server streaming request: %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return fmt.Errorf("Cannot use %v as valid prime factor number\n", number)
	}

	var k int32 = 2

	for number > 1 {
		if number%k == 0 { // if k evenly divides into N
			stream.Send(&calcpb.PrimeResponse{
				PrimeFactor: k,
			})
			number = number / k // divide N by k so that we have the rest of the number left.
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) Average(stream calcpb.Calc_AverageServer) error {
	var count int32 = 0
	var sum int32 = 0
	var average float64 = 0.0
	for {
		req, err := stream.Recv()
		fmt.Printf("Received Average Request with number: %v\n", req.GetNumber())
		if err != nil {
			if err == io.EOF {
				// client stopped sending requests
				average = float64(sum) / float64(count)
				res := calcpb.AverageResponse{
					Average: average,
				}
				stream.SendAndClose(&res)
				break
			} else {
				log.Printf("Error receiving client request: %v\n", err)
				return fmt.Errorf("Error receiving client request: %v\n", err)
			}
		}
		sum += req.GetNumber()
		count++
	}

	return nil
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
