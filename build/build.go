//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/sh"
)

// builds the owner microservice
func Owner(ctx context.Context) error {
	return sh.Run("docker", "build", "-f", "owner.Dockerfile", "..")
}
