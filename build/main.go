package main

import (
	"context"
	"log"

	"dagger.io/dagger"
)

func main() {
	client, err := dagger.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
}
