package main

import (
	"cacher/pkg/cacher"
	"cacher/pkg/client"
	"flag"
	"log"
)

func main() {
	mode := flag.String("mode", "server", "mode = (client, server)")
	flag.Parse()

	if *mode == "server" {
		log.Println("Start in server mode...")
		var s cacher.Server
		s.Run()
	}

	if *mode == "client" {
		log.Println("Start in client mode...")
		var c client.Client
		c.Run()
	}
}
