package asar

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Pack compresses the directory specified by path into an ASAR archive at the destination dest.
func Pack(src string, asar string) error {
	headerJson := map[string]interface{}{
		"files": map[string]interface{}{},
	}
	dir, err := filepath.Abs(src)
	if err != nil {
		return err
	}

	offset := uint64(0)
	files, err := walkDir(dir, headerJson["files"].(map[string]interface{}), &offset)
	if err != nil {
		return err
	}

	header, err := json.Marshal(headerJson)
	if err != nil {
		return err
	}

	jsonSize := len(header)
	// aligns the given size to a multiple of 4.
	size := jsonSize + (4-(jsonSize%4))%4

	header = append(make([]byte, 16), header...)
	header = append(header, make([]byte, size-jsonSize)...)

	writeu32(header[0:4], 4)
	writeu32(header[4:8], uint32(8+size))
	writeu32(header[8:12], uint32(4+size))
	writeu32(header[12:16], uint32(jsonSize))

	if err := os.WriteFile(asar, header, 0664); err != nil {
		return err
	}

	asarFile, err := os.OpenFile(asar, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer asarFile.Close()

	// appends the content of the file specified by filename to the destination file.
	for _, filename := range files {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		_, _ = io.Copy(asarFile, file)
		file.Close()
	}

	return nil
}

// Unpack extracts the contents of the ASAR archive specified by archive into the directory dest.
func Unpack(asar string, dst string) error {
	file, err := os.Open(asar)
	if err != nil {
		return err
	}
	defer file.Close()

	headerSize, jsonMap, err := readHeader(file)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	return iterateEntries(jsonMap, func(val map[string]interface{}, path string) error {
		if offsetStr, ok := val["offset"].(string); ok {
			return extractFile(file, dst, path, offsetStr, headerSize, val)
		} else {
			return os.MkdirAll(filepath.Join(dst, path), 0755)
		}
	})
}
