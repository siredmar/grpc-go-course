package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/siredmar/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to remote grpc server: ", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	fmt.Println("connected")
	// doUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetingRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Edmar",
			LastName:  "Wollnikel",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error rpc greet: %v", err)
	}
	fmt.Println(res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Edmar",
			LastName:  "Wollnikel",
		},
	}

	res, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error rpc greet: %v", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			// server stopped sending
			break
		}
		if err != nil {
			log.Println("error receiving stream: ", err)
		}
		fmt.Println(msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []greetpb.LongGreetRequest{
		greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Edmar",
			},
		},
		greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Wollnik",
			},
		},
		greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Wollmar",
			},
		},
		greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ehrenfried",
			},
		},
	}

	client, err := c.LongGreet(context.Background())
	if err != nil {
		log.Printf("Cannot perform LongGreet: %v", err)
	}
	for _, req := range requests {
		client.Send(&req)
	}

	res, err := client.CloseAndRecv()
	if err != nil {
		log.Printf("Error receiving server response: %v", err)
	}
	fmt.Println(res.GetResult())
}
