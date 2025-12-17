package catalog

import (
	"crypto/sha256"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"fmt"
)

type FileMeta struct {
	Name string
	Size int64
}

const chunkSize = 4096

func FilesCatalog(folder string) ([]FileMeta, error){
	var files []FileMeta

	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		files = append(files, FileMeta{
			Name: d.Name(),
			Size: info.Size(),
		})
		return nil
	})

	return files, err
}

func ChunkHashes(filePath string) ([]string, error){
	var hashes []string

	file, err := os.Open(filePath)
	if err != nil{
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()
	index := 0
	for {
		if fileSize - int64(index*chunkSize) < 0{
			break
		}
		toHash := chunkSize
		if fileSize - int64(index*chunkSize) < int64(chunkSize) {
			toHash = int(fileSize - int64(index*chunkSize))
		}

		buf := make([]byte, toHash)
		n, err := file.Read(buf)

		if err == io.EOF{
			break
		}
		if err != nil{
			return nil, err
		}
		// for _, b := range buf[:n] {
		// 	fmt.Printf("%02x ", b)
		// }
		// fmt.Printf("Raw bytes: %#v\n", buf[:n])
		// fmt.Println()
		// fmt.Println("N:", n)

		// fmt.Println("bufLen", len(buf))
		hash := sha256.Sum256(buf[:n])
		// fmt.Println("hash:", hash)
		hashes = append(hashes, fmt.Sprintf("%x", hash[:]))
		index++
	}
	// fmt.Println("hashes", hashes)

	return hashes, nil
}