package permbits

import (
	"os"
)

const (
	// UserRead is read permissions for the user.
	UserRead os.FileMode = 0o400
	// UserWrite is write permissions for the user.
	UserWrite os.FileMode = 0o200
	// UserExecute is execute permissions for the user.
	UserExecute os.FileMode = 0o100
	// UserAll is all permissions for the user.
	UserAll os.FileMode = UserRead + UserWrite + UserExecute
)

const (
	// UserReadWrite is a helper combination of read and write permissions for the user.
	UserReadWrite os.FileMode = 0o600
	// UserReadWriteExecute is a helper combination of read, write and execute permissions for the user.
	UserReadWriteExecute os.FileMode = 0o700
)

const (
	// GroupRead is read permissions for the group.
	GroupRead os.FileMode = 0o040
	// GroupWrite is write permissions for the group.
	GroupWrite os.FileMode = 0o020
	// GroupExecute is execute permissions for the group.
	GroupExecute os.FileMode = 0o010
	// GroupAll is all permissions for the group.
	GroupAll os.FileMode = GroupRead + GroupWrite + GroupExecute
)

const (
	// GroupReadWrite is a helper combination of read and write permissions for the group.
	GroupReadWrite os.FileMode = 0o060
	// GroupReadWriteExecute is a helper combination of read, write and execute permissions for the group.
	GroupReadWriteExecute os.FileMode = 0o070
)

const (
	// OtherRead is read permissions for others.
	OtherRead os.FileMode = 0o004
	// OtherWrite is write permissions for others.
	OtherWrite os.FileMode = 0o002
	// OtherExecute is execute permissions for others.
	OtherExecute os.FileMode = 0o001
	// OtherAll is all permissions for others.
	OtherAll os.FileMode = OtherRead + OtherWrite + OtherExecute
)

const (
	// OtherReadWrite is a helper combination of read and write permissions for the others.
	OtherReadWrite os.FileMode = 0o006
	// OtherReadWriteExecute is a helper combination of read, write and execute permissions for the others.
	OtherReadWriteExecute os.FileMode = 0o007
)

// Is compares a supplied os.FileMode and returns true if it contains a reference mode.
func Is(mode os.FileMode, is os.FileMode) bool {
	return (mode & is) == is
}
