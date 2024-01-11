package main

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

func compress(data interface{}) ([]byte, error) {
	var compressed bytes.Buffer

	// Encoding the data using gob
	encoder := gob.NewEncoder(&compressed)
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to encode data: %w", err)
	}

	// Compressing the encoded data
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write(compressed.Bytes()); err != nil {
		return nil, fmt.Errorf("failed to write compressed data: %w", err)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zlib writer: %w", err)
	}

	return compressed.Bytes(), nil
}

func decompress(compressedData []byte) (map[string]interface{}, error) {
	r, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, fmt.Errorf("failed to create zlib reader: %w", err)
	}
	defer r.Close()

	// Decompressing the data
	decompressed, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read decompressed data: %w", err)
	}

	// Decoding the data to map[string]interface{}
	var decodedData map[string]interface{}
	decoder := gob.NewDecoder(bytes.NewReader(decompressed))
	if err := decoder.Decode(&decodedData); err != nil {
		return nil, fmt.Errorf("failed to decode data: %w", err)
	}

	return decodedData, nil
}

func main() {
	// Example map[string]interface{}
	data := map[string]interface{}{
		"Id":         1,
		"IsExist":    true,
		"Name":       "John Doe",
		"CreatedAt":  "2022-01-15T15:04:05Z", // Assuming a string representation of time
		"FloatValue": 123.45,
	}

	// Compress data
	compressedData, err := compress(data)
	if err != nil {
		fmt.Printf("Compression error: %v\n", err)
		return
	}

	// Decompress data
	decompressedData, err := decompress(compressedData)
	if err != nil {
		fmt.Printf("Decompression error: %v\n", err)
		return
	}

	// Use the decompressedData as a map
	fmt.Println("Decompressed data:", decompressedData)
}
