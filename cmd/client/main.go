package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/rakhiazfa/go-custom-tcp/config"
)

func messageReceiver(wg *sync.WaitGroup, connection net.Conn) {
	defer wg.Done()

	buffer := make([]byte, 4096)
	for {
		n, err := connection.Read(buffer)
		if err != nil {
			log.Fatal("Failed to receiving message: ", err)
			break
		}

		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Println(message)
	}
}

func messageWriter(wg *sync.WaitGroup, connection net.Conn) {
	defer wg.Done()

	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Failed to reading message: ", err)
			break
		}

		connection.Write([]byte(message))
	}
}

func main() {
	config.LoadConfig()

	addr := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	connection, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Failed to dial connection: ", err)
	}
	defer connection.Close()

	var wg sync.WaitGroup

	wg.Add(1)
	go messageReceiver(&wg, connection)
	wg.Add(1)
	go messageWriter(&wg, connection)

	wg.Wait()
}
