package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/PKURio/quic-go"
	"github.com/PKURio/quic-go/storage"
	"github.com/PKURio/quic-go/utils"
	"io"
)

const (
	addr           = "localhost:4242"
	message        = "foo\x00bar"
	MaxQuicPktSize = 1370
	SendPktSize    = 40
	RcvPktSize     = 1052
	targetFID      = "00000001f5413a6c6142fa779ab00ec51c4c7726"
)

var (
	crc *storage.Item
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

	contentBuf := make([]byte, 2*1024*1024)
	contentRcvNum := 0
	go func(contentBuf []byte, contentRcvNum *int) {
		for {
			rcvBuf := make([]byte, RcvPktSize)
			nLen, err := io.ReadFull(stream, rcvBuf)
			if err != nil {
				fmt.Println("Client receive error: ", err)
			}
			fmt.Printf("Client receive %d bytes\n", nLen)
			content := utils.S2CUnmarshal(rcvBuf)
			copy(contentBuf[int(content.PktIdx)*1024:], content.Payload[:])
			*contentRcvNum++
		}
	}(contentBuf, &contentRcvNum)

	PktIdx := [8]int16{-1, -1, -1, -1, -1, -1, -1, -1} // -1 means empty request
	for i := 0; i < 2*1024; i += 8 {
		for j := 0; j < 8; j++ {
			PktIdx[j] = int16(i + j)
		}
		sendBuf := make([]byte, RcvPktSize)
		sendBuf = utils.C2SMarshal(targetFID, 0, PktIdx[:])

		fmt.Printf("id:%d Client write size: %d\n", i/8, len(sendBuf))
		stream.Write(sendBuf)
		//time.Sleep(100 * time.Millisecond)
	}

	crc := loadCRC()
	for {
		if contentRcvNum == 2*1024 {
			if storage.VerifyBlockCRC(0, contentBuf, crc.CRCArray) {
				fmt.Println("pass crc exam")
			} else {
				fmt.Println("fail to pass crc exam")
			}
			return nil
		}
	}
}

func main() {
	//trace.Start(os.Stderr)
	//defer trace.Stop()

	err := client()
	if err != nil {
		fmt.Println("err: ", err)
	}
}

func loadCRC() (crc *storage.Item) {
	crc, err := storage.LoadCRCFromFile(targetFID)
	if err != nil {
		fmt.Println("failed to load crc from file")
	}
	return
}
