package main

import (
	"cake/config"
	"cake/pkg/route"
	"fmt"
	"log"
)

var server = &route.Server{}

func main() {
	err := config.Load()

	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	server.Initialize()
	server.Listen(config.GetAppConfig().Server.Port)
}
