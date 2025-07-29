// client/main.go
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	socketPath := "/csi/csi.sock"

	// Wait to make sure server is up
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.Write([]byte("Hello from client"))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println("📥 Client received:", string(buf[:n]))
}
