package pp

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
	_, _ = conn.Write([]byte("HELLO\n"))
	// Wait for WELCOME
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if message != "WELCOME\n" {
		fmt.Println("Unexpected response:", message)
		return
	}

	fmt.Println("here")
	err = SendCatalog(conn, "shared")
	if err != nil{
		fmt.Println("Error:", err)
		return
	}

	// Send the file
	// err = SendFile(conn, "ruben-mavarez-4b0WjAX1h64-unsplash.jpg") // change path as needed
	// if err != nil {
	// 	fmt.Println("File send error:", err)
	// }

	// Request a file
	filename := "hello.txt"
	request := fmt.Sprintf("REQUEST:%s\n", filename)
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Request send error:", err)
		return
	}

	// Receive the file
	err = ReceiveFile(conn)
	if err != nil {
		fmt.Println("Receive error:", err)
	}
}