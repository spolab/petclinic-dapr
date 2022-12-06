package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"dagger.io/dagger"
	"github.com/spolab/petclinic/build/task"
)

func main() {
	//
	// Create a context
	//
	ctx := context.Background()
	//
	// Setup Dagger connection
	//
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()
	//
	// Start the build
	//
	id, err := client.Container().
		From("golang:1.19").
		WithExec(task.ApkInstall("go")).
		WithWorkdir("/src").
		Export(ctx, "spolab/petclinic-owner:latest")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
