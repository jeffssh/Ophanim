package main

import (
	"context"
	"log"
	"time"

	"github.com/jeffssh/Ophanim/message"
	"google.golang.org/grpc"
)

const (
	network = "tcp"
	address = "127.0.0.1:9999"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := message.NewReportMessageServiceClient(conn)

	// Contact the server and print out its response.
	mt := message.ReportMessage_Informational
	conf := message.ReportMessage_Certain
	m := "Hello from Sanity, I called the ReportMessage gRPC over tcp at " + address
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.SendReportMessage(ctx, &message.ReportMessage{ReportMessageType: mt, Confidence: conf, Description: m})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
