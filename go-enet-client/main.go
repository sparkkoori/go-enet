package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/sparkkoori/go-enet/enet"
)

import "log"

import "time"

var serverAddress = flag.String("client_address", "localhost:9998", "The address the server is listening on.")

//var client_address = flag.String("server_address","localhost:9997","The address the client will listen on.")

func main() {
	err := enet.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer enet.Deinitialize()

	host, err := enet.CreateHost(nil, 1, 2, 0, 0)
	if err != nil {
		log.Fatalf("Failed to create host: '%s'", err)
	}
	defer host.Destroy()

	serverAddr, err := net.ResolveUDPAddr("udp", *serverAddress)
	if err != nil {
		log.Fatalf("Invalid udp address: '%s'", err)
	}

	server, err := host.Connect(serverAddr, 2, 42)
	if err != nil {
		log.Fatal(err)
	}

	event, err := host.Service(5 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	if event == nil || event.EventType != enet.EventTypeConnect {
		log.Fatal("Failed to connect to server", event)
	}

	for {
		fmt.Print("send> ")
		var input = "junk"
		fmt.Scan(&input)

		if input == "quit" {
			server.Disconnect(42)
			if event, _ = host.Service(3 * time.Second); event != nil {
				switch event.EventType {
				case enet.EventTypeDisconnect:
					fmt.Printf("disconnected")
					return
				default:
					fmt.Println("discarding events during disconnect", event)
				}
			}
			server.DisconnectNow(42)
			return
		}

		payload := []byte(input)
		server.Send(0, payload, enet.FlagReliable)

		for event, err = host.Service(time.Second); true; {
			if err != nil {
				log.Fatal(err)
			}

			if event == nil {
				break
			}

			switch event.EventType {
			case enet.EventTypeConnect:
				fmt.Printf("Connection made, %v\n", event.Data)
			case enet.EventTypeDisconnect:
				fmt.Printf("Disconnection, %v\n", event.Data)
				return
			case enet.EventTypeReceive:
				msg := string(event.Packet)
				fmt.Printf("message: %v\n", msg)
			default:
				log.Fatal("unkown event", event)
			}

			break
		}
	}
}
