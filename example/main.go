package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/PKURio/quic-go"
	"io"
	"io/ioutil"
	"math/big"
	_ "net/http/pprof"
	"os"
)

const (
	addr          = "localhost:4242"
	MaxQuicPktSize = 1370
	message       = "abc\x00abcabcabcabcabcabc\x00abc"
)

// ReadFile 读取一个磁盘文件
func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Open %s error!\n", path)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Read file error!\n")
	}

	return content, nil
}

func serverStart() error {
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
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
			bufRcv := make([]byte, 10)
			bufSend := make([]byte, MaxQuicPktSize)
			copy(bufSend, message)
			for {
				nLen, err := io.ReadFull(stream, bufRcv)
				if err != nil {
					fmt.Println("Server receive error: ", err)
				}
				//fmt.Printf("Server receive '%d' bytes, content: '%s'\n", nLen, bufRcv)
				fmt.Println("Server receive ", nLen, " content: ", bufRcv)
				stream.Write(bufSend)
				fmt.Printf("Server send content: '%s'\n", bufSend)
			}
		}()
	}
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
	err := serverStart()
	if err != nil {
		fmt.Println("err: ", err)
	}
}
