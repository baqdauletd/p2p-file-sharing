package server

import (
	"bufio"
	"fmt"
	"net"
)

func Server(port string){
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil{
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

    for {
        // Accept incoming connections
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }

        // Handle client connection in a goroutine
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()


    // Read and process data from the client
    // ...


	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from peer:", err)
		return
	}
	fmt.Printf("Received from peer: %s", message)
    // Write data back to the client
    // ...
	if message == "HELLO\n" {
		conn.Write([]byte("WELCOME\n"))
	} else {
		conn.Write([]byte("UNKNOWN COMMAND\n"))
	}
}