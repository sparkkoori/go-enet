package main

import (
	"flag"
	"log"
	"time"

	"github.com/sparkkoori/go-enet/enet"
)

import "net"

var address = flag.String("address", "localhost:9998", "The address the server will listen on.")

func main() {
	err := enet.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	defer enet.Deinitialize()

	addr, err := net.ResolveUDPAddr("udp", *address)
	if err != nil {
		log.Fatalf("Invalid udp address: '%s'", err)
	}

	host, err := enet.CreateHost(addr, 2, 2, 0, 0)
	if err != nil {
		log.Fatalf("Failed to create host: '%s'", err)
	}
	defer host.Destroy()
	log.Printf("Start on: '%s'", addr)

	for {
		event, err := host.Service(3 * time.Second)
		if err != nil {
			log.Fatal(err)
		}

		if event == nil {
			continue
		}

		switch event.EventType {
		case enet.EventTypeConnect:
			log.Println("new connection: ", event.Data)
		case enet.EventTypeDisconnect:
			log.Println("disconnection: ", event.Data)
		case enet.EventTypeReceive:
			msg := string(event.Packet)
			log.Println("received: ", msg, event.Packet)
			switch msg {
			case "stop":
				event.Peer.Disconnect(42)
			case "stopall":
				return
			case "die":
				return
			default:
				event.Peer.Send(0, event.Packet, enet.FlagReliable)
			}
		}
	}
}
