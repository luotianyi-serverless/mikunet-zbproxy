package bufio

import (
	"bufio"
	"net"
	"testing"
)

func TestCopyConn(t *testing.T) {
	aa, ab := net.Pipe()
	ba, bb := net.Pipe()
	defer aa.Close()
	defer bb.Close()
	go CopyConn(ab, ba)

	aa.Write([]byte("hello world\r\n"))
	bScanner := bufio.NewScanner(bb)
	if !bScanner.Scan() {
		t.Fail()
	}
	if bScanner.Text() != "hello world" {
		t.Fail()
	}

	bb.Write([]byte("nice to meet you\r\n"))
	aScanner := bufio.NewScanner(aa)
	if !aScanner.Scan() {
		t.Fail()
	}
	if aScanner.Text() != "nice to meet you" {
		t.Fail()
	}
}
