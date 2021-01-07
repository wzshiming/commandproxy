package commandproxy

import (
	"context"
	"net"
	"reflect"
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

func TestDialProxyCommand_Format(t *testing.T) {
	type args struct {
		network string
		address string
	}
	tests := []struct {
		name string
		p    DialProxyCommand
		args args
		want []string
	}{
		{
			p: DialProxyCommand{"nc", "%h", "%p"},
			args: args{
				network: "tcp",
				address: ":1",
			},
			want: []string{"nc", "0.0.0.0", "1"},
		},
		{
			p: DialProxyCommand{"nc", "%h", "%p"},
			args: args{
				network: "tcp",
				address: "[::]:1",
			},
			want: []string{"nc", "::", "1"},
		},
		{
			p: DialProxyCommand{"nc", "-U", "%h"},
			args: args{
				network: "unix",
				address: "./x.socks",
			},
			want: []string{"nc", "-U", "./x.socks"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Format(tt.args.network, tt.args.address); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}
