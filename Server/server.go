package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/quic-go/quic-go"
)

func main() {
	listener, err := quic.ListenAddr(":4242", generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("QUIC server listening on :4242")

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(*conn)
	}
}

func handleConnection(conn quic.Conn) {
	defer conn.CloseWithError(0, "closing connection")

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Println("stream open error:", err)
		return
	}
	defer stream.Close()

	videoFile := "3.mp4" // replace with your video path
	file, err := os.Open(videoFile)
	if err != nil {
		log.Println("failed to open video:", err)
		return
	}
	defer file.Close()

	buf := make([]byte, 64*1024) // 64KB chunks
	for {
		n, err := file.Read(buf)
		if err != nil {
			break
		}
		_, err = stream.Write(buf[:n])
		if err != nil {
			log.Println("stream write error:", err)
			return
		}
		time.Sleep(10 * time.Millisecond) // control streaming rate
	}

	fmt.Println("Video stream finished")
}

func generateTLSConfig() *tls.Config {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, _ := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	keyPEM := tls.Certificate{Certificate: [][]byte{certDER}, PrivateKey: key}
	return &tls.Config{
		Certificates: []tls.Certificate{keyPEM},
		NextProtos:   []string{"my-protocol"},
	}
}
