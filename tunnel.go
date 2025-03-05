package commandproxy

import (
	"context"
	"io"
)

func Tunnel(ctx context.Context, c1, c2 io.ReadWriteCloser, buf1, buf2 []byte) error {
	errCh := make(chan error, 2)
	go func() {
		_, err := io.CopyBuffer(c1, c2, buf1)
		errCh <- err
	}()
	go func() {
		_, err := io.CopyBuffer(c2, c1, buf2)
		errCh <- err
	}()
	defer func() {
		_ = c1.Close()
		_ = c2.Close()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
