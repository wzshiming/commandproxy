package commandproxy

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestProxyCommandUnit(t *testing.T) {
	s := httptest.NewServer(nil)
	defer s.Close()

	cli := s.Client()
	cli.Transport = &http.Transport{
		DialContext: func(ctx context.Context, _ string, address string) (net.Conn, error) {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			conn1, conn2 := net.Pipe()
			go func() {
				err := Tunnel(ctx, conn, conn1)
				if err != nil {
					t.Error(err)
					return
				}
			}()
			return conn2, nil
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	resp, err := cli.Get(s.URL)
	if err != nil {
		t.Fatal(err)
		return
	}
	resp.Body.Close()
}

func TestProxyCommand(t *testing.T) {
	s := httptest.NewServer(nil)
	defer s.Close()

	cli := s.Client()
	cli.Transport = &http.Transport{
		DialContext: func(ctx context.Context, _ string, address string) (net.Conn, error) {
			proxy := ProxyCommand(ctx, "go", "run", "./cmd/commandproxy", s.Listener.Addr().String())
			proxy.Stderr = os.Stderr
			return proxy.Stdio()
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	resp, err := cli.Get(s.URL)
	if err != nil {
		t.Fatal(err)
		return
	}
	resp.Body.Close()
}

func TestProxyCommandExit(t *testing.T) {
	s := httptest.NewServer(nil)
	defer s.Close()

	cli := s.Client()
	cli.Transport = &http.Transport{
		DialContext: func(ctx context.Context, _ string, address string) (net.Conn, error) {
			proxy := ProxyCommand(ctx, "go", "run", "./cmd/commandproxy")
			proxy.Stderr = os.Stderr
			return proxy.Stdio()
		},
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	_, err := cli.Get(s.URL)
	if err == nil {
		t.Fail()
	}

}

func TestProxyCommandNotFoundCmd(t *testing.T) {
	proxy := ProxyCommand(context.Background(), "./notfound")
	proxy.Stderr = os.Stderr
	_, err := proxy.Stdio()
	if err == nil {
		t.Fail()
	}
}