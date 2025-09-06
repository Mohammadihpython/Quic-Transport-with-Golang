# QUIC Video Streaming Project

![QUIC Logo](https://upload.wikimedia.org/wikipedia/commons/3/37/QUIC_logo.svg)

This project demonstrates a **QUIC server-client implementation in Go**, with the ability to **stream video over QUIC**. It uses [`quic-go`](https://github.com/quic-go/quic-go) to provide fast, reliable transport over **UDP**, leveraging QUIC’s features like multiplexed streams, low latency, and improved congestion control.

---

## Overview of QUIC

### 1. QUIC on UDP
- QUIC is a **transport protocol built on UDP**, designed to replace TCP for faster and more reliable connections.
- It provides:
  - **Multiplexed streams**: Multiple independent streams per connection without head-of-line blocking.
  - **0-RTT connection establishment**: Faster handshake.
  - **Built-in encryption**: QUIC is always encrypted using TLS 1.3.
  - **Flow control and congestion control**: Like TCP but more flexible.

**Diagram: QUIC vs TCP over UDP**

![QUIC Architecture](https://user-images.githubusercontent.com/placeholder/quic-diagram.png)


- Each **QUIC packet** is carried over UDP (~1200 bytes by default).
- QUIC splits large streams into frames and manages retransmissions itself.

---

## Project Structure


### Server
- Listens on a QUIC port (default `:4242`).
- Accepts incoming QUIC connections.
- Opens a **bidirectional stream** to send a video file in chunks.
- Handles TLS using a **self-signed certificate** (for testing only).

### Client
- Connects to the server over QUIC.
- Accepts the stream and saves the received bytes into a video file.
- Can later play the video with VLC or MPV.

---

## How to Run

### 1. Server

```bash
cd server
go run main.go

QUIC server listening on :4242

cd client
go run main.go
Video received: received_video.mp4
stream, _ := conn.OpenStreamSync(context.Background())
file, _ := os.Open("video.mp4")
buf := make([]byte, 64*1024)
for {
    n, err := file.Read(buf)
    if n > 0 {
        stream.Write(buf[:n])
    }
    if err != nil { break }
}
stream, _ := conn.AcceptStream(context.Background())
outFile, _ := os.Create("received_video.mp4")
buf := make([]byte, 64*1024)
for {
    n, err := stream.Read(buf)
    if n > 0 {
        outFile.Write(buf[:n])
    }
    if err != nil { break }
}
Flow 
[Client]                 [Server]
   |                        |
   |-- QUIC Handshake ----->|
   |<-- TLS 1.3 Setup ------|
   |                        |
   |-- Open Stream -------->|
   |                        |
   |<-- Send Video Chunks --|
   |                        |
   |-- Receive & Save ----->|
References

QUIC protocol (IETF)
