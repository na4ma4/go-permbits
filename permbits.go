package permbits

import (
	"math/bits"
	"os"
)

const (
	// UserRead is read permissions for the user.
	UserRead os.FileMode = 0o400
	// UserWrite is write permissions for the user.
	UserWrite os.FileMode = 0o200
	// UserExecute is execute permissions for the user.
	UserExecute os.FileMode = 0o100
)

const (
	// UserReadWrite is a helper combination of read and write permissions for the user.
	UserReadWrite os.FileMode = UserRead + UserWrite
	// UserReadWriteExecute is a helper combination of read, write and execute permissions for the user.
	UserReadWriteExecute os.FileMode = UserRead + UserWrite + UserExecute
	// UserAll is all permissions for user.
	UserAll os.FileMode = UserReadWriteExecute
)

const (
	// GroupRead is read permissions for the group.
	GroupRead os.FileMode = 0o040
	// GroupWrite is write permissions for the group.
	GroupWrite os.FileMode = 0o020
	// GroupExecute is execute permissions for the group.
	GroupExecute os.FileMode = 0o010
)

const (
	// GroupReadWrite is a helper combination of read and write permissions for the group.
	GroupReadWrite os.FileMode = GroupRead + GroupWrite
	// GroupReadWriteExecute is a helper combination of read, write and execute permissions for the group.
	GroupReadWriteExecute os.FileMode = GroupRead + GroupWrite + GroupExecute
	// GroupAll is all permissions for groups.
	GroupAll os.FileMode = GroupReadWriteExecute
)

const (
	// OtherRead is read permissions for others.
	OtherRead os.FileMode = 0o004
	// OtherWrite is write permissions for others.
	OtherWrite os.FileMode = 0o002
	// OtherExecute is execute permissions for others.
	OtherExecute os.FileMode = 0o001
)

const (
	// OtherReadWrite is a helper combination of read and write permissions for the others.
	OtherReadWrite os.FileMode = OtherRead + OtherWrite
	// OtherReadWriteExecute is a helper combination of read, write and execute permissions for the others.
	OtherReadWriteExecute os.FileMode = OtherRead + OtherWrite + OtherExecute
	// OtherAll is all permissions for others.
	OtherAll os.FileMode = OtherReadWriteExecute
)

// Is compares a supplied os.FileMode and returns true if it contains a reference mode.
func Is(mode os.FileMode, is os.FileMode) bool {
	return (mode & is) == is
}

// Force takes a base mode and then a list of wanted modes, and makes sure the result
// is true for all the wanted modes.
//
// It is more efficient to supply a list of individual modes than a combination mode.
func Force(mode os.FileMode, want ...os.FileMode) os.FileMode {
	for _, item := range want {
		if item&(item-1) != 0 {
			// multiple bits set
			for i := range bits.Len(uint(item)) {
				if v := item & (1 << i); v != 0 {
					if !Is(mode, v) {
						mode += v
					}
				}
			}

			continue
		}

		// single bit set, so just check it and add it if needed.
		if !Is(mode, item) {
			mode += item
		}
	}

	return mode
}
