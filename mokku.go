package main

import (
	"fmt"
	"mokku/api"
	"mokku/constants"
	"sync"
)

func main() {
	fmt.Println("")
	var servers = &sync.WaitGroup{}
	var settings = api.ParseSettings()
	var storage = api.NewStorage()

	// run server/client
	servers.Add(1)
	go func() {
		switch {
		case settings.Type_of == constants.TYPE_CLIENT:
			api.NewClient(settings).Start()
		case settings.Type_of == constants.TYPE_SERVER:
			api.NewServer(settings, storage).Start()
		}
		servers.Done()
	}()

	// run http server
	if settings.Type_of == constants.TYPE_SERVER {
		servers.Add(1)
		go func() {
			api.NewWeb(settings, storage).Start()
			servers.Done()
		}()
	}

	servers.Wait()
	fmt.Println("done")
}
