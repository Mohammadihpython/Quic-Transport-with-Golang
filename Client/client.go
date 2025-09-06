package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/quic-go/quic-go"
)

func main() {
	addr := "127.0.0.1:4242"

	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"my-protocol"},
	}

	conn, err := quic.DialAddr(context.Background(), addr, tlsConf, nil)
	if err != nil {
		log.Fatal("failed to dial:", err)
	}
	defer conn.CloseWithError(0, "client closing")

	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	outFile := "received_video.mp4"
	file, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf := make([]byte, 64*1024)
	for {
		n, err := stream.Read(buf)
		if n > 0 {
			file.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	fmt.Println("Video received:", outFile)
}
