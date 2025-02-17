package compression

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"Empty", []byte{}},
		{"Short", []byte("hello world")},
		{"Long", bytes.Repeat([]byte("The quick brown fox jumps over the lazy dog. "), 100)},
		{"RepeatedBytes", bytes.Repeat([]byte{0xAA}, 1024)},
		{"AlternatingBytes", bytes.Repeat([]byte{0xAA, 0x55}, 512)},
		{
			"RandomData",
			func() []byte {
				data := make([]byte, 1024)
				rand.Read(data)
				return data
			}(),
		},
		{"MostlyZeros", func() []byte {
			data := make([]byte, 1024)
			data[500] = 1
			return data
		}()},
		{"MostlyOnes", func() []byte {
			data := make([]byte, 1024)
			for i := range data {
				data[i] = 1
			}
			data[500] = 0
			return data
		}()},
		{"Mixed", []byte("This is a test with some repeated bytes and some random data.  1111111IIIIIIIIIllllllllllllllllllllll1177777777777777777777")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := Compress(tt.data)
			if err != nil {
				t.Fatal(err)
			}

			decompressed, err := Decompress(compressed)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(tt.data, decompressed) {
				t.Errorf("Data mismatch for %s: original length %d, decompressed length %d", tt.name, len(tt.data), len(decompressed))
			}
		})
	}
}

func TestCompressionRatio(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"Empty", []byte{}},
		{"Short", []byte("hello world")},
		{"Long", bytes.Repeat([]byte("I must not fear. Fear is the mind-killer. Fear is the little-death that brings total obliteration. I will face my fear. I will permit it to pass over me and through me..."), 92)},
		{"RepeatedBytes", bytes.Repeat([]byte{0xAA}, 1024)},
		{"AlternatingBytes", bytes.Repeat([]byte{0xAA, 0x55}, 512)},
		{"RandomData", func() []byte {
			data := make([]byte, 1024)
			rand.Read(data)
			return data
		}()},
		{"MostlyZeros", func() []byte {
			data := make([]byte, 1024)
			data[500] = 1
			return data
		}()},
		{"MostlyOnes", func() []byte {
			data := make([]byte, 1024)
			for i := range data {
				data[i] = 1
			}
			data[500] = 0
			return data
		}()},
		{"Mixed", []byte("This is a test with some repeated bytes and some random data.  1111111IIIIIIIIIllllllllllllllllllllll1177777777777777777777")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := Compress(tt.data)
			if err != nil {
				t.Fatal(err)
			}

			ratio := float64(len(compressed)) / float64(len(tt.data))
			t.Logf("%s: Original size: %d, Compressed size: %d, Compression ratio: %.2f", tt.name, len(tt.data), len(compressed), ratio)

			switch tt.name {
			case "RepeatedBytes":
				if ratio > 0.05 {
					t.Errorf("RepeatedBytes: Compression ratio too high: %.2f (expected < 0.05)", ratio)
				}
			case "AlternatingBytes":
				if ratio > 0.6 {
					t.Errorf("AlternatingBytes: Compression ratio too high: %.2f (expected < 0.6)", ratio)
				}
			case "RandomData":
				if ratio > 0.95 {
					t.Errorf("RandomData: Compression ratio too high: %.2f (expected < 0.95)", ratio)
				}
			case "MostlyZeros", "MostlyOnes":
				if ratio > 0.05 {
					t.Errorf("%s: Compression ratio too high: %.2f (expected < 0.05)", tt.name, ratio)
				}
			case "Empty":
				if ratio != 0 {
					t.Errorf("Empty: Compression ratio should be 0, got %.2f", ratio)
				}
			case "Short":
				if ratio > 1.2 {
					t.Errorf("Short: Compression ratio too high: %.2f (expected < 1.2)", ratio)
				}
			default:
				if ratio > 1.5 {
					t.Errorf("%s: Compression ratio too high: %.2f (expected < 1.5)", tt.name, ratio)
				}
			}
		})
	}
}

func TestSpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"UTF8BOM", []byte("\xEF\xBB\xBFYou cannot go against the nature of a place without strengthening that nature.")},
		{"Arabic", []byte("الْحَمْدُ لِلَّهِ رَبِّ الْعَالَمِينَ")},
		{"Japanese", []byte("いろはにほへと ちりぬるを")},
		{"Emoji", []byte("🐶 🐱 🐭 🐹 🐰 🦊 🐻 🐼 🐨 🐯 🦁 🐮 🐷 🐸 🐵\n🐔 🐧 🐦 🦆 🦅 🦉 🦇 🐺 🐗 🐴 🦄 🐝 🪱 🐛 🦋\n🐌 🐞 🐜 🪰 🪲 🪳 🦟 🦗 🕷 🦂 🐢 🐍 🦎 🦖\n🦕 🐙 🦑 🦐 🦞 🦀 🐡 🐠 🐟 🐬 🐳 🐋 🦈 🐊 🐅\n🐆 🦓 🦍 🦧 🦣 🐘 🦛 🦏 🐪 🐫 🦒 🦘 🦬 🐃 🐂\n🐄 🐎 🐖 🐏 🐑 🦙 🐐 🦌 🐕 🐩 🦮 🐕🦺 🐈 🐈⬛✨")},
		{"Zalgo", []byte("T̴̸̛̛̛̛͊̅̐̅͛̈́́̍̈́̈̀̇̉̓̉̈́̔̋͑̇̾̀͆̇̌͗́͋̌͌̈́̓̿͛̌͗̈́̅̑͋̓͋͛̈́̇̆̏̊̑̇̅̽̈́̌̓̀͑͗̊̾̊͂̅́̈́̓̉͆̎͛͗̈́̋̉̇͆̅̊̆͛̈́͂̐̈́̏͛͌͂̐̈́͌͐͛͐̇̆́̋͂̏̀͛͋͌̐̈́̇͗̏̇́̿̎̽͋̅̈́̔̈́̆̇͗̅̈́̽̾̏̉͗͛̔͂̊̌̅̈́̊̈́͂̓̌̆̆̾̕̕̕̕̕̕̚͝͠͝͝͝͝͠͝͝͝͝͝͝͠͝")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compressed, err := Compress(tt.data)
			if err != nil {
				t.Fatalf("Compression failed: %v", err)
			}

			decompressed, err := Decompress(compressed)
			if err != nil {
				t.Fatalf("Decompression failed: %v", err)
			}

			if !bytes.Equal(tt.data, decompressed) {
				t.Errorf("Data mismatch for %s", tt.name)
			}

			ratio := float64(len(compressed)) / float64(len(tt.data))
			t.Logf("%s: Original size: %d, Compressed size: %d, Compression ratio: %.2f", tt.name, len(tt.data), len(compressed), ratio)
		})
	}
}
