package commandproxy

import (
	"context"
	"net"
	"os/exec"
)

func ProxyCommand(ctx context.Context, name string, arg ...string) *Proxy {
	return &Proxy{
		Cmd: exec.CommandContext(ctx, name, arg...),
	}
}

type Proxy struct {
	*exec.Cmd
}

func (p *Proxy) Stdio() (net.Conn, error) {
	cmd := p.Cmd
	lp, err := exec.LookPath(p.Path)
	if err != nil {
		return nil, err
	}
	conn1, conn2 := net.Pipe()
	cmd.Path = lp
	cmd.Stdin = conn1
	cmd.Stdout = conn1
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		cmd.Process.Wait()
		conn1.Close()
	}()

	return conn2, nil
}
