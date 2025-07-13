package pp

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	if message != "WELCOME\n" {
		fmt.Println("Unexpected response:", message)
		return
	}

	catalog, err := ReceiveCatalog(reader)
	if err != nil{
		fmt.Println("Error:", err)
	}

	fmt.Println("Received Catalog:")
	for _, file := range catalog{
		fmt.Println("Name: "+file.Name)
	}

	// fmt.Println("here")
	// err = SendCatalog(conn, "shared")
	// if err != nil{
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// send the file
	// err = SendFile(conn, "ruben-mavarez-4b0WjAX1h64-unsplash.jpg") // change path as needed
	// if err != nil {
	// 	fmt.Println("File send error:", err)
	// }

	// request a file
	// filename := "hello.txt"
	fmt.Print("Enter the filename to request: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	filename := scanner.Text()

	request := fmt.Sprintf("REQUEST:%s\n", filename)
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Request send error:", err)
		return
	}

	// receive the file
	err = ReceiveFile(conn)
	if err != nil {
		fmt.Println("Receive error:", err)
		return
	}
}