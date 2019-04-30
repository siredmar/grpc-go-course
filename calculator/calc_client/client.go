package main

import (
	"context"
	"fmt"
	"log"

	"github.com/siredmar/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Println("Could not connect to server: ", err)
	}
	defer cc.Close()

	con := calcpb.NewCalcClient(cc)

	req := &calcpb.CalcRequest{
		First:  10,
		Second: 3,
	}
	res, err := con.Sum(context.Background(), req)
	if err != nil {
		log.Println("Error receving Sum: ", err)
	}
	fmt.Println("Received response from server: ", res.Sum)
}
