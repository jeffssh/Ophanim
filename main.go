package main

import (
	"log"
	"time"

	"github.com/jeffssh/Ophanim/module"
)

func main() {
	modules := module.LoadAllModules("./module/modules/")
	for _, m := range modules {
		log.Printf("Loaded module:\n%+v\n", m)
		err := m.Start()
		if err != nil {
			log.Printf("Error when starting module %s, command %s: %v\n", m.Name, m.Command, err)
		}
	}

	time.Sleep(5 * time.Second)

	for _, m := range modules {
		err := m.Stop()
		if err != nil {
			log.Printf("Error when stopping module %s command %s: %v\n", m.Name, m.Command, err)
		}
	}

}
