package permbits

import (
	"os"
)

const (
	// UserRead is read permissions for the user.
	UserRead os.FileMode = 0400
	// UserWrite is write permissions for the user.
	UserWrite os.FileMode = 0200
	// UserExecute is execute permissions for the user.
	UserExecute os.FileMode = 0100
	// UserAll is all permissions for the user.
	UserAll os.FileMode = UserRead + UserWrite + UserExecute
)

const (
	// GroupRead is read permissions for the group.
	GroupRead os.FileMode = 0040
	// GroupWrite is write permissions for the group.
	GroupWrite os.FileMode = 0020
	// GroupExecute is execute permissions for the group.
	GroupExecute os.FileMode = 0010
	// GroupAll is all permissions for the group.
	GroupAll os.FileMode = GroupRead + GroupWrite + GroupExecute
)

const (
	// OtherRead is read permissions for others.
	OtherRead os.FileMode = 0004
	// OtherWrite is write permissions for others.
	OtherWrite os.FileMode = 0002
	// OtherExecute is execute permissions for others.
	OtherExecute os.FileMode = 0001
	// OtherAll is all permissions for others.
	OtherAll os.FileMode = OtherRead + OtherWrite + OtherExecute
)

func Is(mode os.FileMode, is os.FileMode) bool {
	return (mode & is) == is
}
