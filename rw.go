package asar

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

const MAX_SIZE uint64 = 1<<32 - 1

// readu32 reads a 32-bit unsigned integer from a byte slice in little-endian format.
func readu32(buffer []byte) uint32 {
	return binary.LittleEndian.Uint32(buffer[:4])
}

// writeu32 writes a 32-bit unsigned integer to a byte slice in little-endian format.
func writeu32(buffer []byte, value uint32) {
	binary.LittleEndian.PutUint32(buffer[:4], value)
}

// readHeader reads and parses the header from the given file.
func readHeader(reader *os.File) (uint32, map[string]interface{}, error) {
	headerBuffer := make([]byte, 16)
	if _, err := reader.Read(headerBuffer); err != nil {
		return 0, nil, err
	}

	headerSize := readu32(headerBuffer[4:8])
	jsonSize := readu32(headerBuffer[12:16])

	jsonBuffer := make([]byte, jsonSize)
	if _, err := reader.Read(jsonBuffer); err != nil {
		return 0, nil, err
	}

	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonBuffer, &jsonMap); err != nil {
		return 0, nil, err
	}

	return headerSize + 8, jsonMap, nil
}

// iterateEntries recursively iterates through the entries in the JSON map and calls the callback for each entry.
func iterateEntries(jsonMap map[string]interface{}, callback func(map[string]interface{}, string) error) error {
	var helper func(map[string]interface{}, string) error
	helper = func(current map[string]interface{}, path string) error {
		if err := callback(current, path); err != nil {
			return err
		}
		if files, ok := current["files"].(map[string]interface{}); ok {
			for key, val := range files {
				if err := helper(val.(map[string]interface{}), filepath.Join(path, key)); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if files, ok := jsonMap["files"].(map[string]interface{}); ok {
		for key, val := range files {
			if err := helper(val.(map[string]interface{}), key); err != nil {
				return err
			}
		}
	}
	return nil
}

// walkDir walks through the directory tree rooted at dir, updating the JSON map and collecting file paths.
func walkDir(dir string, jsonMap map[string]interface{}, offset *uint64) ([]string, error) {
	var files []string
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() {
			subMap := map[string]interface{}{"files": map[string]interface{}{}}
			jsonMap[name] = subMap
			subFiles, err := walkDir(filepath.Join(dir, name), subMap["files"].(map[string]interface{}), offset)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			fs, _ := entry.Info()
			size := fs.Size()
			if uint64(size) > MAX_SIZE {
				return nil, fmt.Errorf("file %s (%f GB) is above the maximum possible size of %f GB", name, float64(size)/1e9, float64(MAX_SIZE)/1e9)
			}
			jsonMap[name] = map[string]interface{}{
				"offset": strconv.FormatUint(*offset, 10),
				"size":   size,
			}
			*offset += uint64(size)
			files = append(files, filepath.Join(dir, name))
		}
	}
	return files, nil
}

// extractFile extracts a single file from the ASAR archive to the specified destination path.
func extractFile(file *os.File, dst, path, offsetStr string, headerSize uint32, val map[string]interface{}) error {
	offset, err := strconv.ParseUint(offsetStr, 10, 64)
	if err != nil {
		return err
	}
	size := uint64(val["size"].(float64))
	if _, err := file.Seek(int64(headerSize)+int64(offset), io.SeekStart); err != nil {
		return err
	}
	buffer := make([]byte, size)
	if _, err := file.Read(buffer); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dst, path), buffer, 0644)
}
