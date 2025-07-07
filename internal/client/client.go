package client

import (
	"bufio"
	"fmt"
	"net"
	"p2p-file-sharing/internal/transfer"
)

func Client(host, port string) {
    // Connect to the server
    conn, err := net.Dial("tcp", host+":"+port)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer conn.Close()

    // Send data to the server
    // ...
	_, _ = conn.Write([]byte("HELLO\n"))
	// Wait for WELCOME
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if message != "WELCOME\n" {
		fmt.Println("Unexpected response:", message)
		return
	}

	// Send the file
	err = transfer.SendFile(conn, "ruben-mavarez-4b0WjAX1h64-unsplash.jpg") // change path as needed
	if err != nil {
		fmt.Println("File send error:", err)
	}
}