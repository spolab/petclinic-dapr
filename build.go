//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/sh"
)

// builds the owner microservice
func Owner(ctx context.Context) error {
	//
	// Execute the Docker build
	//
	return sh.Run("docker", "build", "-f", "build/owner.Dockerfile", "-t", "spolab/petclinic-owner:latest", ".")
}

func Toolbox(ctx context.Context) error {
	//
	// Execute the Docker build
	//
	return sh.Run("docker", "build", "-f", "build/toolbox.Dockerfile", "-t", "spolab/toolbox:latest", "build")
}
