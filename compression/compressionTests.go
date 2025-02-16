package compression_test

import (
	"bytes"
	"crypto/rand"
	"github.com/idmcarvalho/DataCompression/compression"
	"strings"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"Empty", []byte{}},
		{"RandomData", func() []byte {
			data := make([]byte, 1024)
			rand.Read(data)
			return data
		}()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := compression.Compress(tt.data)
			if err != nil {
				t.Fatal(err)
			}

			decompressed, err := compression.Decompress(compressed)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(tt.data, decompressed) {
				t.Error("Data mismatch")
			}
		})
	}
}
