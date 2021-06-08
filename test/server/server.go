package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"github.com/PKURio/quic-go"
	"github.com/PKURio/quic-go/log"
	"github.com/PKURio/quic-go/node"
	"github.com/PKURio/quic-go/storage"
	"github.com/PKURio/quic-go/utils"
	"io"
	"math/big"
	"net"
	_ "net/http/pprof"
	"strconv"
)

const (
	addr              = "localhost:4242"
	message           = "foo\x00bar"
	MaxQuicPktSize    = 1370
	ClientSendPktSize = 40
	ClientRcvPktSize  = 1052
	ServerSendPktSize = 1052
	ServerRcvPktSize  = 40
	targetFID         = "00000001f5413a6c6142fa779ab00ec51c4c7726"
)

var (
	data [6][]byte
)

func server() error {
	listener, err := quic.ListenAddrEarly(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	for {
		sess, err := listener.Accept(context.Background())
		if err != nil {
			return err
		}

		// process
		go func() {
			stream, _ := sess.AcceptStream(context.Background())
			rcvBuf := make([]byte, ServerRcvPktSize)
			sendBuf := make([]byte, ServerSendPktSize)
			for {
				nLen, err := io.ReadFull(stream, rcvBuf)
				if err != nil {
					fmt.Println("Server receive error: ", err)
				}
				fmt.Printf("Server receive %d content %s\n", nLen, rcvBuf)
				content := utils.C2SUnmarshal(rcvBuf)
				fmt.Println("FID", hex.EncodeToString(content.FID[:]), "FileIdx", content.FileIdx, "PktIdx", content.PktIdx)

				for _, idx := range content.PktIdx {
					if idx < 0 {
						continue
					}
					fid := hex.EncodeToString(content.FID[:])
					// TODO: modify payload size for tail packet
					sendBuf = utils.S2CMarshal(fid, 0, idx, int16(1024), data[0][int(idx)*1024:(int(idx)+1)*1024])

					stream.Write(sendBuf)
					fmt.Printf("PktIdx:%d Server send size:%d\n", idx, len(sendBuf))
				}
			}
		}()
	}
}

// External interface to start server
func ServerStart(conn net.PacketConn, path string) error {
	log.GetLogger().Println("ServerStart.")
	node.Conn = conn
	storage.Path = path
	err := server()
	return err
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}

func main() {
	node.Conn = &net.UDPConn{}
	loadData()

	err := server()
	if err != nil {
		fmt.Println("err: ", err)
	}
}

func loadData() {
	for i := 0; i < 6; i++ {
		data[i], _ = storage.ReadFile("../../data/" + targetFID + "_" + strconv.Itoa(i) + storage.DataFileExtension)
		fmt.Printf("data[%d] size: %d\n", i, len(data[i]))
	}
}
