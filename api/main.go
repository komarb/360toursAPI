package main

import (
	"./config"
	"./server"
)

func main() {
	config := config.GetConfig()
	app := &server.App{}
	app.Init(config)
	app.Run(":8080")
}