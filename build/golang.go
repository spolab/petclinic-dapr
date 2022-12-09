// Standard recipes for the magefiles
package build

import (
	"github.com/magefile/mage/sh"
)

func GoCompile(source string, target string) error {
	return sh.Run("go", "build", "-o", target, "-buildvcs=true", source)
}
