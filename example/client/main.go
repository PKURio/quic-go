package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/PKURio/quic-go"
	"github.com/PKURio/quic-go/node"
	"io"
	"net"
	"time"
)

const (
	addr          = "localhost:4242"
	message       = "foo\x00bar"
	MaxQuicPktSize = 1370
)

func client() error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
	session, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return err
	}

	stream, err := session.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	var cnt uint8 = 0
	go func(cntPtr *uint8) {
		bufRcv := make([]byte, MaxQuicPktSize)
		for {
			nLen, err := io.ReadFull(stream, bufRcv)
			if err != nil {
				fmt.Println("Client receive error: ", err)
			}
			fmt.Printf("Client receive '%d' bytes, content: '%s'\n", nLen, bufRcv)
			//fmt.Println("Client receive ", nLen, " content: ", bufRcv)
			*cntPtr++
		}
	}(&cnt)

	bufSend := make([]byte, 10)
	copy(bufSend, message)
	for {
		stream.Write(bufSend)
		fmt.Printf("Client write content: '%s'\n", bufSend)
		time.Sleep(500 * time.Millisecond)
	}
}

func clientStart(conn net.PacketConn) error {
	node.Conn = conn



	return nil
}


func main() {
	err := client()
	if err != nil {
		fmt.Println("err: ", err)
	}
}
