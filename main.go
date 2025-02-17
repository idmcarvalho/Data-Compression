package main

import (
	"DataCompression/compression"
	"bytes"
	"fmt"
	"log"
)

func main() {
	original := bytes.Repeat([]byte("The highest function of ecology is the understanding of consequences."), 100)

	compressed, err := compression.Compress(original)
	if err != nil {
		log.Fatal(err)
	}

	decompressed, err := compression.Decompress(compressed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original size: %d bytes\n", len(original))
	fmt.Printf("Compressed size: %d bytes (%.1f%% of original)\n",
		len(compressed),
		float64(len(compressed))/float64(len(original))*100)
	fmt.Println("Original:", string(original))
	fmt.Println("Decompressed:", string(decompressed))

	if len(compressed) >= len(original) {
		fmt.Println("\n⚠️ Compression didn't reduce size - try with:")
		fmt.Println("- Larger input data")
		fmt.Println("- More repetitive content")
	}
}
