# permbits

Golang File Permission Bit Operators.

[![CI](https://github.com/na4ma4/permbits/workflows/CI/badge.svg)](https://github.com/na4ma4/permbits/actions?query=workflow%3ACI)
[![GoDoc](https://godoc.org/github.com/na4ma4/permbits/src/jwt?status.svg)](https://godoc.org/github.com/na4ma4/permbits)

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

## Constants

| Mode | Constant |
| ---- | -------- |
| u+a  | permbits.UserAll |
| u+r  | permbits.UserRead |
| u+w  | permbits.UserWrite |
| u+x  | permbits.UserExecute |
| g+a  | permbits.GroupAll |
| g+r  | permbits.GroupRead |
| g+w  | permbits.GroupWrite |
| g+x  | permbits.GroupExecute |
| o+a  | permbits.OtherAll |
| o+r  | permbits.OtherRead |
| o+w  | permbits.OtherWrite |
| o+x  | permbits.OtherExecute |
