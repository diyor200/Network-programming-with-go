package main

import (
	"bufio"
	"net"
	"reflect"
	"testing"
)

const payload = "Diyorbek Abdulaxatov"

var expected = []string{"Diyorbek", "Abdulaxatov"}

func TestScanner(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte(payload))
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanWords)

	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}
	
	if !reflect.DeepEqual(expected, words) {
		t.Fatal("inaccurate scanned word list")
	}
	t.Logf("Scanned words: %#v", words)
}
