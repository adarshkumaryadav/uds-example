package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	socketPath := "/csi/csi.sock"

	// Remove if already exists
	_ = os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("🔌 Server listening on", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("📩 Server received:", string(buf[:n]))
	conn.Write([]byte("Hello from server via UDS"))
}
