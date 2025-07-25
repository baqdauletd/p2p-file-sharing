package pp

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
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

	// send metadata
	_, _ = writer.WriteString("FILENAME:" + stat.Name() + "\n")
	_, _ = writer.WriteString("SIZE:" + fmt.Sprintf("%d\n", stat.Size()))
	writer.Flush()

	hashes, err := GenerateChunkHashes(filepath)
	if err != nil{
		return err
	}

	// fmt.Println("Hashes in SendFile:", hashes)

	_, _ = writer.WriteString("HASHCOUNT:" + fmt.Sprintf("%d\n", len(hashes)))
	for _, h := range hashes {
		_, _ = writer.WriteString(h + "\n")
	}
	writer.Flush()


	// send file in chunks
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

	// read filename
	filenameLine, err := reader.ReadString('\n')
	if err != nil{
		return err
	}
	filename := strings.TrimPrefix(strings.TrimSpace(filenameLine), "FILENAME:")

	// read size
	sizeLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	sizeStr := strings.TrimPrefix(strings.TrimSpace(sizeLine), "SIZE:")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		fmt.Println("here1")
		return err
	}

	hashCountLine, err := reader.ReadString('\n')
	if err != nil{
		fmt.Println("here2")
		return err
	}
	hashCountStr := strings.TrimPrefix(strings.TrimSpace(hashCountLine), "HASHCOUNT:")
	hashCount, err := strconv.Atoi(hashCountStr)
	if err != nil{
		fmt.Println("here3")
		return err
	}
	expectedHashes := make([]string, 0, hashCount)
	for i := 0; i < hashCount; i++ {
		h, _ := reader.ReadString('\n')
		expectedHashes = append(expectedHashes, strings.TrimSpace(h))
	}

	targetDir := "received-files"
	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		fmt.Print("Error creating directory:")
		return err
	}
	fmt.Println("Created Directory")

	fullPath := filepath.Join(targetDir, "received_" + filename)


	outFile, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("here4")
		return err
	}
	defer outFile.Close()

	// copy chunks
	written := 0
	buf := make([]byte, chunkSize)
	chunkIndex := 0

	for written < size {
		toRead := chunkSize
		if size-written < chunkSize {
			toRead = size - written
		}
		n, err := io.ReadFull(reader, buf[:toRead])
		if err != nil && err != io.EOF {
			fmt.Println("here5")
			return err
		}

		// verify chunk hash
		hash := sha256.Sum256(buf[:n])
		hashStr := fmt.Sprintf("%x", hash[:])

		if hashStr != expectedHashes[chunkIndex] {
			fmt.Println("here6")
			return fmt.Errorf("chunk %d hash mismatch", chunkIndex)
		}

		outFile.Write(buf[:n])
		written += n
		chunkIndex++
	}

	fmt.Println("File received:", "received_"+filename)
	return nil
}

func SendCatalog(conn net.Conn, folder string) error {
	files, err := BuildFileCatalog(folder)
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(conn)
	writer.WriteString("CATALOG\n")
	for _, f := range files {
		writer.WriteString(fmt.Sprintf("FILENAME:%s;SIZE:%d\n", f.Name, f.Size))
	}
	writer.WriteString("ENDCATALOG\n")
	return writer.Flush()
}

func ReceiveCatalog(reader *bufio.Reader) ([]FileMeta, error) {
	var catalog []FileMeta

	firstLine, _ := reader.ReadString('\n')
	if strings.TrimSpace(firstLine) != "CATALOG" {
		return nil, fmt.Errorf("missing CATALOG header")
	}
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return catalog, err
		}
		line = strings.TrimSpace(line)
		if line == "ENDCATALOG" {
			break
		}
		if strings.HasPrefix(line, "FILENAME:") {
			parts := strings.Split(line, ";")
			if len(parts) < 2 {
				continue
			}
			name := strings.TrimPrefix(parts[0], "FILENAME:")
			sizeStr := strings.TrimPrefix(parts[1], "SIZE:")
			size, _ := strconv.ParseInt(sizeStr, 10, 64)
			catalog = append(catalog, FileMeta{
				Name: name,
				Size: size,
			})
		}
	}
	return catalog, nil
}
