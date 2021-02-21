package main

import (
	"github.com/scottapp/misuchiru-server/server"
)

func main() {
	server := server.NewServer()
	server.SetupRouter()
	go server.HttpServer.ListenAndServe()
	server.GracefulShutdown(3000)
}
