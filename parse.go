package permbits

import (
	"errors"
	"fmt"
	"os"
)

type parseMode int

const (
	parseOther = 0
	parseGroup = 3
	parseUser  = 6
)

type parseApply int

const (
	parseApplyNone parseApply = iota
	parseApplyAdd
	parseApplySet
	parseApplySub
)

// ErrModeStringSyntax is returned when there is a syntax error in the mode string.
var ErrModeStringSyntax = errors.New("syntax error in mode string")

// MustString is a wrapper around FromString that will panic on invalid syntax.
func MustString(perms string) os.FileMode {
	mode, err := FromString(perms)
	if err != nil {
		panic(err)
	}

	return mode
}

// FromString takes a subset of the available symbolic modes and returns a os.FileMode
// that is comprised of them or 0 and an error if the input is invalid.
//
// Supported Modes
//   References: ugoa
//   Operators: +-=
//   Modes: rwx
//
// Special modes are not supported.
//
// Because the mode starts at 0 the `-` operator does not remove permissions unless they
// were added earlier in the mode string.
//nolint:cyclop,funlen // it's a long function.
func FromString(perms string) (os.FileMode, error) {
	o := 0
	mode := []parseMode{}
	apply := parseApplyNone

	for i := 0; i < len(perms); i++ {
		t := perms[i]
		switch t {
		case 'u':
			mode = append(mode, parseUser)
		case 'g':
			mode = append(mode, parseGroup)
		case 'o':
			mode = append(mode, parseOther)
		case 'a':
			mode = append(mode, parseOther, parseGroup, parseUser)
		case ',':
			mode = []parseMode{}
			apply = parseApplyNone
		case '+':
			apply = parseApplyAdd
		case '-':
			apply = parseApplySub
		case '=':
			apply = parseApplySet
		case 'r', 'w', 'x':
			v := 0

			switch t {
			case 'r':
				v = 4
			case 'w':
				v = 2
			case 'x':
				v = 1
			}

			for _, pv := range mode {
				iv := v << pv

				switch apply {
				case parseApplyNone:
					return os.FileMode(o), fmt.Errorf("%w: %s", ErrModeStringSyntax, perms)
				case parseApplyAdd, parseApplySet:
					if (o & iv) != iv {
						o += iv
					}
				case parseApplySub:
					if (o & iv) == iv {
						o -= iv
					}
				}
			}
		default:
			return os.FileMode(o), fmt.Errorf("%w: %s", ErrModeStringSyntax, perms)
		}
	}

	return os.FileMode(o), nil
}
