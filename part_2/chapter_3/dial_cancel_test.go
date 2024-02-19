package ch3

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

func TestDialContextCance(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	syncc := make(chan struct{})

	go func() {
		defer func() {
			syncc <- struct{}{}
		}()

		var d net.Dialer
		d.Control = func(_, _ string, c syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}

		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Log(err)
			return
		}

		conn.Close()
		t.Error("connection did not time out")
	}()

	cancel()
	<-syncc

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual %v", ctx.Err())
	}
}
