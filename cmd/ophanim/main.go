package main

import (
	"log"
	"net"
	"time"

	"github.com/jeffssh/Ophanim/message"
	"github.com/jeffssh/Ophanim/module"
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
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	log.Printf("Hosting message gRPC server using %s at %s", network, address)
	modules := module.LoadAllModules("./module/modules/")
	for _, m := range modules {
		sep := "======================================="
		log.Printf("Loaded module:\n%s\n%+v\n%s", sep, m, sep)
		err := m.Start()
		if err != nil {
			log.Printf("Error when starting module %s, command %s: %v", m.Name, m.Command, err)
		}
	}

	time.Sleep(10 * time.Second)

	for _, m := range modules {
		err := m.Stop()
		if err != nil {
			log.Printf("Error when stopping module %s command %s: %v", m.Name, m.Command, err)
		}
	}

}
