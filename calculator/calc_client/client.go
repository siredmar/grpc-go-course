package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"

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

	primeRequest := calcpb.PrimeRequest{
		Number: 120,
	}

	primestream, err := con.Prime(context.Background(), &primeRequest)
	for {
		res, err := primestream.Recv()
		if err == io.EOF {
			// server stopped sending
			break
		}
		if err != nil {
			log.Println("Error receiving from server")
		}
		fmt.Println(res.GetPrimeFactor())
	}

	averagestream, err := con.Average(context.Background())

	for i := 0; i < 10; i++ {
		randomnumber := rand.Int31n(10)
		req := calcpb.AverageRequest{
			Number: randomnumber,
		}
		fmt.Printf("Sending average request with number: %v\n", randomnumber)

		averagestream.Send(&req)
	}
	averageres, err := averagestream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(averageres.GetAverage())
}
