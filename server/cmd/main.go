package main

import (
	"DataImpactProject/server/app/server"
	"log"
)

func main() {
	server, err := server.NewDataImpactServer()
	if err != nil {
		log.Fatal(err)
	}
	initRouter(server.Router)
	//TODO add port as an env var
	server.Router.Run(":8081")
}
