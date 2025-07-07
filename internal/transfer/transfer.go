package transfer

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const chunkSize = 4096

func SendFile(conn net.Conn, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil{
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(conn)

	// Send metadata
	_, _ = writer.WriteString("FILENAME:" + stat.Name() + "\n")
	_, _ = writer.WriteString("SIZE:" + fmt.Sprintf("%d\n", stat.Size()))
	writer.Flush()

	// Send file in chunks
	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		conn.Write(buf[:n])
	}

	fmt.Println("File sent:", stat.Name())
	return nil
}

func ReceiveFile(conn net.Conn) error{
	reader := bufio.NewReader(conn)

	// Read FILENAME
	filenameLine, err := reader.ReadString('\n')
	if err != nil{
		return err
	}
	filename := strings.TrimPrefix(strings.TrimSpace(filenameLine), "FILENAME:")

	// Read SIZE
	sizeLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	sizeStr := strings.TrimPrefix(strings.TrimSpace(sizeLine), "SIZE:")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return err
	}

	// Prepare to write
	outFile, err := os.Create("received_" + filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Copy chunks
	written := 0
	buf := make([]byte, chunkSize)

	for written < size {
		toRead := chunkSize
		if size-written < chunkSize {
			toRead = size - written
		}
		n, err := io.ReadFull(reader, buf[:toRead])
		if err != nil && err != io.EOF {
			return err
		}
		outFile.Write(buf[:n])
		written += n
	}

	fmt.Println("File received:", "received_"+filename)
	return nil
}
