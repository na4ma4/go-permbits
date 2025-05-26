package permbits_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/na4ma4/go-permbits"
)

func testSymbolicModes(t *testing.T, modeString string, is uint32) {
	t.Helper()

	mode, err := permbits.FromString(modeString)
	if err != nil {
		t.Errorf("[%s] error occurred: %s", modeString, err)
		return
	}

	if uint32(mode) != is {
		t.Errorf("[%s] expected %o, go %o", modeString, mode, is)
	}
}

func TestParseSymbolicModes(t *testing.T) {
	list := map[string]uint32{
		"a+r":         0o444,
		"a-x":         0o000,
		"a+rx":        0o555,
		"u=rw,g=r,o=": 0o640,
		"u+w,go-w":    0o200,
		"a+w,go-w":    0o200,
		"ug=rw":       0o660,
	}
	for k, v := range list {
		testSymbolicModes(t, k, v)
	}
}

func TestParseSymbolicMode1(t *testing.T) {
	modeString := "u+r,g+r,o+r"
	mode, err := permbits.FromString(modeString)
	if err != nil {
		t.Errorf("[%s] error occurred: %s", modeString, err)
		return
	}

	if uint32(mode) != 0o444 {
		t.Errorf("[%s] expected %o, go 0o444", modeString, mode)
	}

	if !permbits.Is(mode, permbits.UserRead) {
		t.Errorf("[%s] expected u+r to be true", modeString)
	}
	if permbits.Is(mode, permbits.UserWrite) {
		t.Errorf("[%s] expected u+w to be false", modeString)
	}
	if permbits.Is(mode, permbits.UserExecute) {
		t.Errorf("[%s] expected u+x to be false", modeString)
	}

	if !permbits.Is(mode, permbits.GroupRead) {
		t.Errorf("[%s] expected g+r to be true", modeString)
	}
	if permbits.Is(mode, permbits.GroupWrite) {
		t.Errorf("[%s] expected g+w to be false", modeString)
	}
	if permbits.Is(mode, permbits.GroupExecute) {
		t.Errorf("[%s] expected g+x to be false", modeString)
	}

	if !permbits.Is(mode, permbits.OtherRead) {
		t.Errorf("[%s] expected o+r to be true", modeString)
	}
	if permbits.Is(mode, permbits.OtherWrite) {
		t.Errorf("[%s] expected o+w to be false", modeString)
	}
	if permbits.Is(mode, permbits.OtherExecute) {
		t.Errorf("[%s] expected o+x to be false", modeString)
	}
}

func TestParseWeirdValue(t *testing.T) {
	modeString := "u+rg+ro+r"
	_, err := permbits.FromString(modeString)
	if err != nil {
		t.Errorf("[%s] error occurred: %s", modeString, err)
	}
}

func TestParseReturnError(t *testing.T) {
	modeString := "a:r"
	_, err := permbits.FromString(modeString)
	if err == nil {
		t.Errorf("[%s] no error occurred", modeString)
	}
}

func testShellChmod(t *testing.T, modeString string) {
	t.Helper()

	pbMode, err := permbits.FromString(modeString)
	if err != nil {
		t.Errorf("[%s] error occurred: %s", modeString, err)
		return
	}

	chMode, err := shellTest(modeString)
	if err != nil {
		t.Errorf("[%s] error occurred: %s", modeString, err)
		return
	}

	if chMode != pbMode {
		t.Errorf("[%s] shell mode (%o) and permbits mode (%o) are not identical", modeString, chMode, pbMode)
	}
}

func TestShellChangeMode(t *testing.T) {
	list := []string{
		"a+r",
		"a-x",
		"a+rx",
		"u=rw,g=r,o=",
		"u+w,go-w",
		"ug=rw",
		"a=r,g-w",
	}
	for _, v := range list {
		testShellChmod(t, v)
	}
}

func TestFromString_ResolveSymbolic(t *testing.T) {
	tests := []struct {
		name       string
		modeString string
		is         os.FileMode
	}{
		{"a+r", "a+r", 0o444},
		{"a-x", "a-x", 0o000},
		{"a+rx", "a+rx", 0o555},
		{"u=rw,g=r,o=", "u=rw,g=r,o=", 0o640},
		{"u+w,go-w", "u+w,go-w", 0o200},
		{"a+w,go-w", "a+w,go-w", 0o200},
		{"ug=rw", "ug=rw", 0o660},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, err := permbits.FromString(tt.modeString)
			if err != nil {
				t.Errorf("permbits.FromString() error = %v", err)
			}
			if mode != tt.is {
				t.Errorf("permbits.FromString() = 0o%04o, want = 0o%04o", mode, tt.is)
			}
		})
	}
}

func TestFromString_ResolveSymbolicExample(t *testing.T) {
	type tv struct {
		read, write, execute bool
	}
	tests := []struct {
		name               string
		modeString         string
		is                 os.FileMode
		user, group, other tv
	}{
		{"u+r,g+r,o+r", "u+r,g+r,o+r", 0o444, tv{true, false, false}, tv{true, false, false}, tv{true, false, false}},
		{"u+w,g+w,o+w", "u+w,g+w,o+w", 0o222, tv{false, true, false}, tv{false, true, false}, tv{false, true, false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mode, err := permbits.FromString(tt.modeString)
			if err != nil {
				t.Errorf("permbits.FromString() error = %v", err)
			}
			if mode != tt.is {
				t.Errorf("permbits.FromString() = 0o%04o, want = 0o%04o", mode, tt.is)
			}

			pbtest := func(mode, match os.FileMode, expect bool) {
				if v := permbits.Is(mode, match); v != expect {
					t.Errorf(
						"permbits.FromString() = 0o%04o to return %t for 0o%04o, but returned %t",
						mode, expect, match, v,
					)
				}
			}

			pbtest(mode, permbits.UserRead, tt.user.read)
			pbtest(mode, permbits.UserWrite, tt.user.write)
			pbtest(mode, permbits.UserExecute, tt.user.execute)
			pbtest(mode, permbits.GroupRead, tt.group.read)
			pbtest(mode, permbits.GroupWrite, tt.group.write)
			pbtest(mode, permbits.GroupExecute, tt.group.execute)
			pbtest(mode, permbits.OtherRead, tt.other.read)
			pbtest(mode, permbits.OtherWrite, tt.other.write)
			pbtest(mode, permbits.OtherExecute, tt.other.execute)
		})
	}
}

func TestFromString_ValidParseNoError(t *testing.T) {
	tests := []struct {
		expect     fs.FileMode
		modeString string
	}{
		{0o0444, "u+rg+ro+r"},
		{0o0440, "u+rg+r"},
	}
	for _, tt := range tests {
		t.Run(tt.modeString, func(t *testing.T) {
			mode, err := permbits.FromString(tt.modeString)
			if err != nil {
				t.Errorf("permbits.FromString() error = %v", err)
			}

			if mode != tt.expect {
				t.Errorf("permbits.FromString(): got '0o%04o', want '0o%04o'",
					mode,
					tt.expect,
				)
			}
		})
	}
}

func TestFromString_InvalidParseError(t *testing.T) {
	tests := []string{
		"a:r",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			if _, err := permbits.FromString(tt); err == nil {
				t.Errorf("permbits.FromString() should return error for %s", tt)
			}
		})
	}
}
