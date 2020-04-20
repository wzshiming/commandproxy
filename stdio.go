package commandproxy

import (
	"os"
)

var Stdio stdio

type stdio struct{}

func (stdio) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (stdio) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (stdio) Close() (err error) {
	os.Exit(0)
	return nil
}
