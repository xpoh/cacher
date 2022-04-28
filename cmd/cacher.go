package main

import (
	"github.com/xpoh/cacher/pkg/cacher"
	"github.com/xpoh/cacher/pkg/client"
	"time"
)

func main() {
	var s cacher.Server
	go s.Run()

	time.Sleep(time.Second)

	var c client.Client
	c.Run()
}
