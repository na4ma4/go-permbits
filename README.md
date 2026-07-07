# go-permbits

Golang File Permission Bit Operators.

[![CI](https://github.com/na4ma4/go-permbits/actions/workflows/ci.yml/badge.svg)](https://github.com/na4ma4/go-permbits/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/na4ma4/go-permbits?status.svg)](https://godoc.org/github.com/na4ma4/go-permbits)

## Usage

Anywhere you can specify a file mode, you can use the permbits to add together the permissions you want.

(This test works because it's likely the umask is `022` which will mean the first MkdirAll will actually create a directory with 0750 permissions).

```golang
testPath := "./test"
mode := permbits.UserAll+permbits.GroupAll

if err := os.MkdirAll(testPath, mode); err != nil {
    log.Fatal(err)
}

if st, err := os.Stat(testPath); err == nil {
    if !permbits.Is(st.Mode(), mode) {
        log.Printf("Updating mode for %s (was %o)", testPath, st.Mode())
        os.Chmod(testPath, mode)
    } else {
        log.Printf("Test path mode is correct: %s (%o)", testPath, st.Mode())
    }
}
```

```shell
$ go run main.go
2021/07/09 13:12:53 Updating mode for ./test (was 20000000750)
```

## Functions

### `Is(mode, is os.FileMode) bool`

Checks whether `mode` contains all the bits set in `is`.

### `Force(mode os.FileMode, want ...os.FileMode) os.FileMode`

Takes a base mode and ensures all wanted permission bits are set.

```golang
mode := permbits.Force(0o111, permbits.UserReadWriteExecute)
// mode is now 0o711
```

### `FromString(perms string) (os.FileMode, error)`

Parses symbolic mode strings (a subset of `chmod` syntax) into an `os.FileMode`.

Supports references `ugoa`, operators `+-=`, and modes `rwx`.

```golang
mode, err := permbits.FromString("u=rw,g=r,o=")
// mode is 0o640
```

### `MustString(perms string) os.FileMode`

Wrapper around `FromString` that panics on invalid input.

## Constants

### Permission Bits

| Mode | Constant |
| ---- | -------- |
| u+r  | `permbits.UserRead` |
| u+w  | `permbits.UserWrite` |
| u+x  | `permbits.UserExecute` |
| g+r  | `permbits.GroupRead` |
| g+w  | `permbits.GroupWrite` |
| g+x  | `permbits.GroupExecute` |
| o+r  | `permbits.OtherRead` |
| o+w  | `permbits.OtherWrite` |
| o+x  | `permbits.OtherExecute` |

### Helper Combinations

| Mode | Constant |
| ---- | -------- |
| u+r+w    | `permbits.UserReadWrite` |
| u+r+w+x  | `permbits.UserReadWriteExecute` or `permbits.UserAll` |
| g+r+w    | `permbits.GroupReadWrite` |
| g+r+w+x  | `permbits.GroupReadWriteExecute` or `permbits.GroupAll` |
| o+r+w    | `permbits.OtherReadWrite` |
| o+r+w+x  | `permbits.OtherReadWriteExecute` or `permbits.OtherAll` |
