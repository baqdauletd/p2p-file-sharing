package pp

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func GenerateChunkHashes(filePath string) ([]string, error){
	var hashes []string

	file, err := os.Open(filePath)
	if err != nil{
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, chunkSize)
	for {
		n, err := file.Read(buf)
		if err == io.EOF{
			break
		}
		if err != nil{
			return nil, err
		}

		hash := sha256.Sum256(buf[:n])
		hashes = append(hashes, fmt.Sprintf("%x", hash[:]))
	}

	return hashes, nil
}