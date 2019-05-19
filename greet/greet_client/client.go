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
	doUnary(c)
	doServerStreaming(c)
	doClientStreaming(c)
	doBidirectionalStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("------- Do Unary ----------------------")
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
	fmt.Println("------- Do Server Streaming------------")
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
	fmt.Println("------- Do Client Streaming------------")
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

func doBidirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("------- Do Bidirectional Streaming-----")
	client, err := c.ManyGreets(context.Background())
	if err != nil {
		log.Printf("Cannot perform LongGreet: %v", err)
	}
	waitc := make(chan bool)

	go func() {
		requests := []greetpb.ManyGreetsRequest{
			greetpb.ManyGreetsRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Edmar",
				},
			},
			greetpb.ManyGreetsRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Wollnik",
				},
			},
			greetpb.ManyGreetsRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Wollmar",
				},
			},
			greetpb.ManyGreetsRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "Ehrenfried",
				},
			},
		}

		for _, req := range requests {
			err := client.Send(&req)
			if err != nil {
				fmt.Printf("Error sending request: %v\n", err)
			}
		}
		client.CloseSend()
	}()

	go func() {
		for {
			res, err := client.Recv()
			if err != nil {
				if err == io.EOF {
					// server stopped streaming
					break
				} else {
					fmt.Printf("Error receiving: %v\n", err)
				}
			}
			fmt.Println(res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}
