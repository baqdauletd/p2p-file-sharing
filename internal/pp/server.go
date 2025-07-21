package pp

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	// "log"
	"net"
)

func Server(port string){
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil{
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port "+port)

    for {
        // Accept incoming connections
		fmt.Println("waiting for connection")
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error:", err)
            continue
        }
		// log.Println("here2")
        // Handle client connection in a goroutine
		fmt.Println("connected")
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()

	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from peer:", err)
		return
	}
	fmt.Printf("Received from peer: %s\n", message)


	if message == "HELLO\n" {
		conn.Write([]byte("WELCOME\n"))

		// Receive file after handshake
		// err := ReceiveFile(conn)
		// if err != nil {
		// 	fmt.Println("Receive error:", err)
		// }

		err = SendCatalog(conn, "shared")
		if err != nil{
			fmt.Println("Error:", err)
			return
		}

		// fmt.Println("READIING IN SERVER")
		// for{
		// 	msg3, err := reader.ReadString('\n')
		// 	fmt.Println(msg3)
		// 	if err == io.EOF{
		// 		break
		// 	}
		// }
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}
		line = strings.TrimSpace(line)
		// fmt.Printf("message after sending a catalog: %s\n", line)

		if strings.HasPrefix(line, "REQUESTCHUNK:") {
			parts := strings.Split(line, ":")
			if len(parts) != 3 {
				fmt.Println("Malformed chunk request")
				return
			}
			filename := parts[1]
			index, err := strconv.Atoi(parts[2])
			if err != nil {
				fmt.Println("Invalid chunk index")
				return
			}
			err = handleChunkRequest(conn, filename, index)
			if err != nil {
				fmt.Println("Chunk send error:", err)
			}
			// fmt.Println("Ending CHUNK")
			return
		}

		if strings.HasPrefix(line, "REQUEST:") {
			filename := strings.TrimPrefix(line, "REQUEST:")
			filePath := "shared/" + filename

			fmt.Println("Peer requested file:", filename)

			err := SendFile(conn, filePath)
			if err != nil {
				fmt.Println("File send error:", err)
			}
		} else {
			fmt.Println("Unknown command after handshake:", line)
		}
	} else {
		conn.Write([]byte("UNKNOWN COMMAND\n"))
	}
}

func handleChunkRequest(conn net.Conn, filename string, index int) error {
	filePath := "shared/" + filename
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	offset := int64(index * chunkSize)
	buf := make([]byte, chunkSize)

	_, err = file.ReadAt(buf, offset)
	if err != nil && err != io.EOF {
		return err
	}

	_, err = conn.Write(buf)
	return err
}
