package commandproxy

import (
	"context"
	"net"
	"testing"
)

func TestDialProxyCommand(t *testing.T) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}

	command, err := SplitCommand("nc %n %p")
	if err != nil {
		t.Fatal(err)
	}
	dial := DialProxyCommand(command)
	conn, err := dial.DialContext(context.Background(), listener.Addr().Network(), listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}
