//go:build darwin || linux

package permbits_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/na4ma4/go-permbits"
)

//nolint:wrapcheck // test function
func shellTest(mode string) (os.FileMode, error) {
	tmpFile, err := os.CreateTemp("", "permbits-test-*")
	if err != nil {
		return 0, err
	}

	name := tmpFile.Name()
	tmpFile.Close()
	defer func() {
		_ = os.Remove(name)
	}()

	cmd := exec.Command(`chmod`, "ugo=", name)
	if err = cmd.Run(); err != nil {
		return 0, err
	}

	cmd = exec.Command(`chmod`, mode, name)
	if err = cmd.Run(); err != nil {
		return 0, err
	}

	stv, err := os.Stat(name)
	if err != nil {
		return 0, err
	}

	return stv.Mode(), nil
}

func TestFromString_CompareToCommandLine(t *testing.T) {
	tests := []string{
		"a+r",
		"a-x",
		"a+rx",
		"u=rw,g=r,o=",
		"u+w,go-w",
		"ug=rw",
		"a=r,g-w",
	}
	for _, tt := range tests {
		t.Run("Compare to command line: "+tt, func(t *testing.T) {
			pbMode, err := permbits.FromString(tt)
			if err != nil {
				t.Errorf("permbits.FromString() error = %v", err)
			}

			chMode, err := shellTest(tt)
			if err != nil {
				t.Errorf("shellTest error = %v", err)
			}

			if chMode != pbMode {
				t.Errorf("permbits.FromString() = %04o, shellTest() = %04o", pbMode, chMode)
			}
		})
	}
}
