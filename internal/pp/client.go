package pp

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
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
	fmt.Printf("Request line is: %s\n", request)
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Request send error:", err)
		return
	}

	// receive the file
	// err = ReceiveFile(conn)
	// if err != nil {
	// 	fmt.Println("Receive error:", err)
	// 	return
	// }
	// for{
	// 	msg3, err := reader.ReadString('\n')
	// 	fmt.Println(msg3)
	// 	if err == io.EOF{
	// 		break
	// 	}
	// }

	filenameLine, err := reader.ReadString('\n')
	fmt.Printf("fileNameLine: %s\n", filenameLine)
	if err != nil {
		fmt.Println("Failed to read filename:", err)
		return
	}
	fileName := strings.TrimPrefix(strings.TrimSpace(filenameLine), "FILENAME:")
	fmt.Printf("fileName: %s\n", fileName)

	sizeLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read size:", err)
		return
	}
	sizeStr := strings.TrimPrefix(strings.TrimSpace(sizeLine), "SIZE:")
	fileSize, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid size:", err)
		return
	}
	fmt.Printf("fileSize: %d\n", fileSize)

	hashCountLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read hash count:", err)
		return
	}
	hashCountStr := strings.TrimPrefix(strings.TrimSpace(hashCountLine), "HASHCOUNT:")
	hashCount, err := strconv.Atoi(hashCountStr)
	if err != nil {
		fmt.Println("Invalid hash count:", err)
		return
	}
	fmt.Printf("hashCount: %d\n", hashCount)

	// read hashes
	hashes := make([]string, 0, hashCount)
	for i := 0; i < hashCount; i++ {
		h, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read hash:", err)
			return
		}
		hashes = append(hashes, strings.TrimSpace(h))
	}
	fmt.Printf("hashes: %s\n", hashes)

	err = ConcurrentChunks(host, port, fileName, hashes, fileSize, reader)
	if err != nil {
		fmt.Println("Concurrent receive error:", err)
	return
}
}

func ConcurrentChunks(host, port, fileName string, hashes []string, fileSize int64, reader *bufio.Reader) error{
	totalChunks := len(hashes)
	chunkData := make([][]byte, totalChunks)
	errChan := make(chan error, totalChunks)
	var wg sync.WaitGroup

	for a := 0; a < totalChunks; a++{
		wg.Add(1)
		go func(index int){
			defer wg.Done()

			conn, err := net.Dial("tcp", host+":"+port)
			if err != nil{
				errChan <- err
				return
			}
			defer conn.Close()

			conn.Write([]byte("HELLO\n"))
			bufReader := bufio.NewReader(conn)

			// for{
			// 	msg3, err := bufReader.ReadString('\n')
			// 	fmt.Println(msg3)
			// 	if err == io.EOF{
			// 		break
			// 	}
			// }

			msg, err := bufReader.ReadString('\n')
			if err != nil{
				fmt.Println("here11")
				errChan <- err
				return
			}

			fmt.Printf("msg: %s\n", msg)


			if msg != "WELCOME\n"{
				errChan <-fmt.Errorf("unexpected response")
				return
			}

			request := fmt.Sprintf("REQUESTCHUNK:%s:%d\n", fileName, index)
			conn.Write([]byte(request))

			// fmt.Println("IN SECOND TIME")
			// for{
			// 	msg3, err := bufReader.ReadString('\n')
			// 	fmt.Println(msg3)
			// 	if err == io.EOF{
			// 		break
			// 	}
			// }

			// buf := make([]byte, chunkSize)
			toRead := chunkSize
			fmt.Println("fileSize:", fileSize)
			fmt.Println("chunkSize:", chunkSize)
			if fileSize - int64(index*chunkSize) < int64(chunkSize) {
				toRead = int(fileSize - int64(index*chunkSize))
			}
			buf := make([]byte, toRead)
			n, err := io.ReadFull(reader, buf)

			// for _, b := range buf[:n] {
			// 	fmt.Printf("%02x ", b)
			// }
			fmt.Printf("Raw bytes: %#v\n", buf[:n])
			fmt.Println()

			fmt.Println("N:",n)

			if err != nil && err != io.EOF{
				fmt.Println("here12")
				errChan <- fmt.Errorf("chunk %d read error: %w", index, err)
				return
			}


			fmt.Println("bufLen", len(buf))
			hash := sha256.Sum256(buf[:n])
			fmt.Println("hash:", hash)
			hashStr := fmt.Sprintf("%x", hash[:])
			fmt.Printf("hashStr: %s\n", hashStr)
			if hashStr != hashes[index]{
				errChan <- fmt.Errorf("hash mismatch at chunk %d", index)
				return
			}
			chunkData[index] = append(chunkData[index], buf[:n]...)

		}(a)
	}

	wg.Wait()
	close(errChan)

	if len(errChan) > 0{
		for er := range errChan{
			fmt.Println("Error: ", er)
		}
		return fmt.Errorf("some chunks failed")

	}

	outFile, err := os.Create("received_" + fileName)
	if err != nil {
		return err
	}
	defer outFile.Close()

	for _, chunk := range chunkData{
		_, err := outFile.Write(chunk)
		if err != nil{
			return err
		}
	}


	fmt.Println("Download complete.")
	return nil
}