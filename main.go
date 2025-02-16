package main

import (
	"fmt"
	"github.com/idmcarvalho/DataCompression/compression"
	"log"
)

func main() {
	data := []byte("Test data")
	compressed, err := compression.Compress(data)
	if err != nil {
		log.Fatal(err)
	}

	decompressed, err := compression.Decompress(compressed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Original:", string(data))
	fmt.Println("Decompressed:", string(decompressed))
}
