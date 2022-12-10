//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/sh"
)

// build the owner image (spolab/petclinic-owner)
func Owner(ctx context.Context) error {
	return sh.Run("docker", "build", "-f", "owner/docker/Dockerfile", "-t", "spolab/petclinic-owner:latest", ".")
}

// build the toolbox image (spolab/toolbox)
func Toolbox(ctx context.Context) error {
	return sh.Run("docker", "build", "-t", "spolab/toolbox:latest", "toolbox")
}
