package main

import (
	"log"
	"net"

	"github.com/jeffssh/Ophanim/cli"
	"github.com/jeffssh/Ophanim/message"
)

const (
	network = "tcp"
	address = "127.0.0.1:9999"
)

func main() {

	lis, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := message.NewMessageServer()
	log.Printf("Hosting message gRPC server using %s at %s", network, address)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	c := cli.New()
	for {
		c.Prompt()
	}
}

//𓁹𓁹 𓁺 𓁻 𓁼 𓁽 𓁾 𓁿
