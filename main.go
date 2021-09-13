package main

import (
	"log"
	"strconv"
	"tcp-fragmented-server/internal/pkg/client"
	"tcp-fragmented-server/internal/pkg/server"
	"time"
)

func main() {
	address := "0.0.0.0:9999"

	srv := server.New(address)
	go func() {
		log.Fatal(srv.Serve())
	}()
	time.Sleep(100 * time.Millisecond) // it's not for production, so just sleep while server is starting

	clients := make([]*client.Client, 0, 10)

	for i := 0; i < 10; i++ {
		name := "#" + strconv.Itoa(i) // this name for logging purposes. Real client IDs are set automatically by server
		c := client.New(address, name)
		if err := c.Connect(); err != nil {
			log.Fatal(err)
		}
		go func() {
			err := c.Run()
			if err != nil {
				log.Println(err)
			}
		}()
		clients = append(clients, c)
	}

	// Message to tagged client
	err := clients[1].SendMessage("@2, this is message from client #1 to client #2")
	if err != nil {
		log.Fatal(err)
	}

	// Message to all clients
	//err = clients[3].SendMessage("this message to all")
	//if err != nil {
	//	log.Fatal(err)
	//}

	time.Sleep(15 * time.Second)
}
