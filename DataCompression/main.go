package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"log"
)

func main() {
	original := []byte("I must not fear. Fear is the mind-killer. Fear is the little-death that brings total obliteration. I will face my fear. I will permit it to pass over me and through me...")

	compressed, err := compressData(original)
	if err != nil {
		log.Fatalf("Compression failed: %v", err)
	}

	fmt.Printf("Original size: %d bytes\n", len(original))
	fmt.Printf("Compressed size: %d bytes (%.2f%% reduction)\n",
		len(compressed),
		100*(1-float64(len(compressed))/float64(len(original))))

	decompressed, err := decompressData(compressed)
	if err != nil {
		log.Fatalf("Decompression failed: %v", err)
	}

	if !bytes.Equal(original, decompressed) {
		log.Fatal("Decompressed data doesn't match original!")
	}

	fmt.Println("\nCompression/decompression successful!")
}

func compressData(data []byte) ([]byte, error) {
	var buf bytes.Buffer

	zw, err := flate.NewWriterDict(&buf, flate.BestCompression, nil)
	if err != nil {
		return nil, fmt.Errorf("compression error: %w", err)
	}

	if _, err := zw.Write(data); err != nil {
		zw.Close()
		return nil, fmt.Errorf("compression write error: %w", err)
	}

	if err := zw.Close(); err != nil {
		return nil, fmt.Errorf("compression close error: %w", err)
	}

	return buf.Bytes(), nil
}

func decompressData(compressed []byte) ([]byte, error) {
	r := flate.NewReader(bytes.NewReader(compressed))
	defer r.Close()

	var decompressed bytes.Buffer
	if _, err := io.Copy(&decompressed, r); err != nil {
		return nil, fmt.Errorf("decompression error: %w", err)
	}

	if _, err := r.Read(make([]byte, 1)); err != io.EOF {
		return nil, fmt.Errorf("corrupt data: unexpected content after decompression")
	}

	return decompressed.Bytes(), nil
}
