package main

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestExamplePinger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	r, w := io.Pipe()
	done := make(chan struct{})
	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second // initial ping interval

	go func() {
		Pinger(ctx, w, resetTimer)
		close(done)
	}()

	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			fmt.Printf("resetting timer (%s)\n", d)
			resetTimer <- d
		}

		now := time.Now()
		buf := make([]byte, 1024)
		n, err := r.Read(buf)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("received %q (%s)\n", buf[:n], time.Since(now).Round(100*time.Millisecond))
	}

	for i, v := range []int64{0, 200, 300, -1, -1, -1} {
		fmt.Printf("Run %d:\n", i+1)
		receivePing(time.Duration(v)*time.Millisecond, r)
	}

	cancel()
	<-done // ensures the pinger exits after canceling the context
}
