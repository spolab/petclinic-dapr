//go:build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/spolab/petclinic/build"
)

type Build mg.Namespace

// vets the full package
func Vet() error {
	return sh.Run("go", "vet", "./...")
}

// tests the full package
func Test() error {
	return sh.Run("go", "test", "./...")
}

// builds the owner actor
func (Build) Actor() error {
	mg.Deps(Test, Vet)
	return build.GoCompile("cmd/actor/actor.go", "bin/app")
}

// builds the owner user interface
func (Build) UI() error {
	mg.Deps(Test, Vet)
	return build.GoCompile("cmd/ui/ui.go", "bin/app")
}
