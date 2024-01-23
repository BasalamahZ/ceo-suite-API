package main

import (
	"os"
	"sync"

	"github.com/ceo-suite/cmd/ceosuite-api-http/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		os.Exit(server.Run())
		defer wg.Done()
	}()
	wg.Wait()
}
