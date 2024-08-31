package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abdelmounim-dev/redis/pkg/server"
)

func main() {
	s, err := server.NewServer(":6379", 100, time.Second*5)
	if err != nil {
		log.Fatal(err)
	}
	s.Run()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Printf("killing the server :D")
	s.Kill()
}
