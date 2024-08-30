package main

import (
	"log"

	"github.com/abdelmounim-dev/redis/pkg/server"
)

func main() {
	s, err := server.NewServer(6379, 100)
	if err != nil {
		log.Fatal(err)
	}
	s.Run()

}
