package main

import "babo/utility/zlog"

func main() {
	// Initialize the client application
	zlog.Init("client_virtual", false, false)
	defer zlog.Sync()
	client := &Client{}
	client.Start()
}
