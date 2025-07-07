package client

import (
	"bufio"
	"fmt"
	"net"
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
	data := []byte("HELLO\n")
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

    // Read and process data from the server
    // ...
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	fmt.Printf("Response from peer: %s", message)
}