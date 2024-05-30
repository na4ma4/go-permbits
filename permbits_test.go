package permbits_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/na4ma4/go-permbits"
)

const (
	modeAll  = os.FileMode(0o777)
	modeNone = os.FileMode(0o000)
)

func TestIs_CompareSingleModeAll(t *testing.T) {
	tests := []struct {
		name string
		mode os.FileMode
		is   os.FileMode
	}{
		{"Compare 0o777 to UserRead", modeAll, permbits.UserRead},
		{"Compare 0o777 to UserWrite", modeAll, permbits.UserWrite},
		{"Compare 0o777 to UserExecute", modeAll, permbits.UserExecute},
		{"Compare 0o777 to GroupRead", modeAll, permbits.GroupRead},
		{"Compare 0o777 to GroupWrite", modeAll, permbits.GroupWrite},
		{"Compare 0o777 to GroupExecute", modeAll, permbits.GroupExecute},
		{"Compare 0o777 to OtherRead", modeAll, permbits.OtherRead},
		{"Compare 0o777 to OtherWrite", modeAll, permbits.OtherWrite},
		{"Compare 0o777 to OtherExecute", modeAll, permbits.OtherExecute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v := permbits.Is(tt.mode, tt.is); !v {
				t.Errorf("permbits.Is() returned %t for %04o matching %04o", v, tt.mode, tt.is)
			}
		})
	}
}

func TestIs_CompareSingleModeNone(t *testing.T) {
	tests := []struct {
		name string
		mode os.FileMode
		is   os.FileMode
	}{
		{"Compare 0o000 to UserRead", modeNone, permbits.UserRead},
		{"Compare 0o000 to UserWrite", modeNone, permbits.UserWrite},
		{"Compare 0o000 to UserExecute", modeNone, permbits.UserExecute},
		{"Compare 0o000 to GroupRead", modeNone, permbits.GroupRead},
		{"Compare 0o000 to GroupWrite", modeNone, permbits.GroupWrite},
		{"Compare 0o000 to GroupExecute", modeNone, permbits.GroupExecute},
		{"Compare 0o000 to OtherRead", modeNone, permbits.OtherRead},
		{"Compare 0o000 to OtherWrite", modeNone, permbits.OtherWrite},
		{"Compare 0o000 to OtherExecute", modeNone, permbits.OtherExecute},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v := permbits.Is(tt.mode, tt.is); v {
				t.Errorf("permbits.Is() returned %t for %04o matching %04o", v, tt.mode, tt.is)
			}
		})
	}
}

func TestIs_CompareMultipleModes(t *testing.T) {
	tests := []struct {
		name string
		mode os.FileMode
		is   os.FileMode
		want bool
	}{
		{
			"Compare 0o777 to UserAll+GroupAll+OtherAll, expect True",
			0o777,
			permbits.UserAll + permbits.GroupAll + permbits.OtherAll,
			true,
		},
		{
			"Compare 0o775 to UserAll+GroupAll+OtherAll, expect False",
			0o775,
			permbits.UserAll + permbits.GroupAll + permbits.OtherAll,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v := permbits.Is(tt.mode, tt.is); v != tt.want {
				t.Errorf("permbits.Is() returned %t, expected %t for %04o matching %04o", v, tt.want, tt.mode, tt.is)
			}
		})
	}
}

func TestForce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		mode   os.FileMode
		expect os.FileMode
		want   []os.FileMode
	}{
		{
			0o111, 0o711,
			[]os.FileMode{permbits.UserReadWriteExecute},
		},
		{
			0o111, 0o711,
			[]os.FileMode{permbits.UserRead, permbits.UserWrite, permbits.UserExecute},
		},
		{
			0o111, 0o171,
			[]os.FileMode{permbits.GroupRead, permbits.GroupWrite, permbits.GroupExecute},
		},
	}
	for _, tt := range tests {
		t.Run(
			fmt.Sprintf("%04o + %v = %04o", tt.mode, tt.want, tt.expect),
			func(t *testing.T) {
				t.Parallel()

				if v := permbits.Force(tt.mode, tt.want...); v != tt.expect {
					t.Errorf("permbits.Force() got '%04o', want '%04o'", v, tt.expect)
				}
			},
		)
	}
}
