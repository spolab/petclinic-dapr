//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Owner mg.Namespace

// build the owner image (spolab/petclinic-owner-actor)
func (Owner) Actor(ctx context.Context) error {
	return sh.Run("nerdctl", "build", "--namespace", "k8s.io", "--build-arg", "FILE_SRC=owner/cli/actor/actor.go", "-f", "owner/docker/Dockerfile", "-t", "spolab/petclinic-owner-actor:latest", ".")
}

// build the owner image (spolab/petclinic-owner-ui)
func (Owner) Api(ctx context.Context) error {
	return sh.Run("nerdctl", "build", "--namespace", "k8s.io", "--build-arg", "FILE_SRC=owner/cli/api/api.go", "-f", "owner/docker/Dockerfile", "-t", "spolab/petclinic-owner-api:latest", ".")
}

// build the owner image (spolab/petclinic-owner-ui)
func (Owner) Ui(ctx context.Context) error {
	if err := sh.Run("npm", "--prefix", "owner/ui", "run", "build"); err != nil {
		return err
	}
	return sh.Run("nerdctl", "build", "--namespace", "k8s.io", "-f", "owner/ui/Dockerfile", "-t", "spolab/petclinic-owner-ui:latest", "owner/ui")
}

// build all the owner services
func (Owner) All() {
	mg.Deps(Owner.Actor, Owner.Ui)
}

// build the toolbox image (spolab/toolbox)
func Toolbox(ctx context.Context) error {
	return sh.Run("nerdctl", "build", "--namespace", "k8s.io", "-t", "spolab/toolbox:latest", "toolbox")
}

func GitHash(filename string) (hash string, err error) {
	hash, err = sh.Output("git", "rev-parse", "HEAD")
	return
}
