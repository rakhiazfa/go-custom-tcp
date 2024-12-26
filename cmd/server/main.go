package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rakhiazfa/go-custom-tcp/config"
	"github.com/rakhiazfa/go-custom-tcp/core"
)

func main() {
	config.LoadConfig()

	addr := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	server := core.NewServer(addr)

	err := server.Run(func() {
		log.Printf("Server running at %s\n", addr)
	})
	if err != nil {
		log.Fatalln("Failed to run server: ", err)
	}
}
