package polyline_test

import (
	"math/rand"
	"testing"

	"github.com/twpayne/go-polyline"
)

// generateCoords generates n random coordinates in valid lat/lon ranges
func generateCoords(n int, seed int64) [][]float64 {
	r := rand.New(rand.NewSource(seed))
	coords := make([][]float64, n)
	for i := range coords {
		coords[i] = []float64{
			180*r.Float64() - 90,  // lat: -90 to 90
			360*r.Float64() - 180, // lon: -180 to 180
		}
	}
	return coords
}

// generateFlatCoords generates n*2 random flat coordinates
func generateFlatCoords(n int, seed int64) []float64 {
	r := rand.New(rand.NewSource(seed))
	coords := make([]float64, n*2)
	for i := 0; i < n*2; i += 2 {
		coords[i] = 180*r.Float64() - 90
		coords[i+1] = 360*r.Float64() - 180
	}
	return coords
}

// Pre-encode test data for decode benchmarks
var (
	encoded10      []byte
	encoded100     []byte
	encoded1000    []byte
	encoded10000   []byte
	coords10       [][]float64
	coords100      [][]float64
	coords1000     [][]float64
	coords10000    [][]float64
	flatCoords100  []float64
	flatEncoded100 []byte
)

func init() {
	coords10 = generateCoords(10, 42)
	coords100 = generateCoords(100, 42)
	coords1000 = generateCoords(1000, 42)
	coords10000 = generateCoords(10000, 42)

	encoded10 = polyline.EncodeCoords(coords10)
	encoded100 = polyline.EncodeCoords(coords100)
	encoded1000 = polyline.EncodeCoords(coords1000)
	encoded10000 = polyline.EncodeCoords(coords10000)

	flatCoords100 = generateFlatCoords(100, 42)
	codec := polyline.Codec{Dim: 2, Scale: 1e5}
	flatEncoded100, _ = codec.EncodeFlatCoords(nil, flatCoords100)
}

// ============================================================================
// DecodeUint Benchmarks
// ============================================================================

func BenchmarkDecodeUint(b *testing.B) {
	// Small value (1 byte): 0-31
	b.Run("small_1byte", func(b *testing.B) {
		buf := polyline.EncodeUint(nil, 15) // Encodes to 1 byte
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeUint(buf)
		}
	})

	// Medium value (3 bytes): typical coordinate delta
	b.Run("medium_3bytes", func(b *testing.B) {
		buf := polyline.EncodeUint(nil, 50000) // Typical delta, ~3 bytes
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeUint(buf)
		}
	})

	// Large value (6 bytes): large coordinate
	b.Run("large_6bytes", func(b *testing.B) {
		buf := polyline.EncodeUint(nil, 500000000) // Large value, ~6 bytes
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeUint(buf)
		}
	})
}

// ============================================================================
// EncodeUint Benchmarks
// ============================================================================

func BenchmarkEncodeUint(b *testing.B) {
	b.Run("small_1byte", func(b *testing.B) {
		buf := make([]byte, 0, 16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeUint(buf[:0], 15)
		}
	})

	b.Run("medium_3bytes", func(b *testing.B) {
		buf := make([]byte, 0, 16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeUint(buf[:0], 50000)
		}
	})

	b.Run("large_6bytes", func(b *testing.B) {
		buf := make([]byte, 0, 16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeUint(buf[:0], 500000000)
		}
	})
}

// ============================================================================
// DecodeInt Benchmarks
// ============================================================================

func BenchmarkDecodeInt(b *testing.B) {
	b.Run("positive", func(b *testing.B) {
		buf := polyline.EncodeInt(nil, 3850000) // Typical lat*1e5
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeInt(buf)
		}
	})

	b.Run("negative", func(b *testing.B) {
		buf := polyline.EncodeInt(nil, -12020000) // Typical lon*1e5
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeInt(buf)
		}
	})

	b.Run("small_delta", func(b *testing.B) {
		buf := polyline.EncodeInt(nil, 100) // Small delta between coords
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeInt(buf)
		}
	})
}

// ============================================================================
// EncodeInt Benchmarks
// ============================================================================

func BenchmarkEncodeInt(b *testing.B) {
	b.Run("positive", func(b *testing.B) {
		buf := make([]byte, 0, 16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeInt(buf[:0], 3850000)
		}
	})

	b.Run("negative", func(b *testing.B) {
		buf := make([]byte, 0, 16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeInt(buf[:0], -12020000)
		}
	})

	b.Run("small_delta", func(b *testing.B) {
		buf := make([]byte, 0, 16)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeInt(buf[:0], 100)
		}
	})
}

// ============================================================================
// DecodeCoords Benchmarks
// ============================================================================

func BenchmarkDecodeCoords(b *testing.B) {
	b.Run("n=10", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded10)
		}
	})

	b.Run("n=100", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded100)
		}
	})

	b.Run("n=1000", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded1000)
		}
	})

	b.Run("n=10000", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded10000)
		}
	})
}

// ============================================================================
// EncodeCoords Benchmarks
// ============================================================================

func BenchmarkEncodeCoords(b *testing.B) {
	b.Run("n=10", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords10)
		}
	})

	b.Run("n=100", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords100)
		}
	})

	b.Run("n=1000", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords1000)
		}
	})

	b.Run("n=10000", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords10000)
		}
	})
}

// ============================================================================
// DecodeFlatCoords Benchmarks
// ============================================================================

func BenchmarkDecodeFlatCoords(b *testing.B) {
	codec := polyline.Codec{Dim: 2, Scale: 1e5}

	b.Run("n=100", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = codec.DecodeFlatCoords(nil, flatEncoded100)
		}
	})

	b.Run("n=1000", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _, _ = codec.DecodeFlatCoords(nil, encoded1000)
		}
	})
}

// ============================================================================
// EncodeFlatCoords Benchmarks
// ============================================================================

func BenchmarkEncodeFlatCoords(b *testing.B) {
	codec := polyline.Codec{Dim: 2, Scale: 1e5}
	flatCoords1000 := generateFlatCoords(1000, 42)

	b.Run("n=100", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = codec.EncodeFlatCoords(nil, flatCoords100)
		}
	})

	b.Run("n=1000", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = codec.EncodeFlatCoords(nil, flatCoords1000)
		}
	})
}

// ============================================================================
// Codec Scale Variations
// ============================================================================

func BenchmarkDecodeCoords_Scale(b *testing.B) {
	// Encode with different scales
	codec1e5 := polyline.Codec{Dim: 2, Scale: 1e5}
	codec1e6 := polyline.Codec{Dim: 2, Scale: 1e6}
	codec1e7 := polyline.Codec{Dim: 2, Scale: 1e7}

	encoded1e5 := codec1e5.EncodeCoords(nil, coords100)
	encoded1e6 := codec1e6.EncodeCoords(nil, coords100)
	encoded1e7 := codec1e7.EncodeCoords(nil, coords100)

	b.Run("scale=1e5", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = codec1e5.DecodeCoords(encoded1e5)
		}
	})

	b.Run("scale=1e6", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = codec1e6.DecodeCoords(encoded1e6)
		}
	})

	b.Run("scale=1e7", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, _ = codec1e7.DecodeCoords(encoded1e7)
		}
	})
}

// ============================================================================
// Memory Allocation Benchmarks
// ============================================================================

func BenchmarkDecodeCoords_Allocs(b *testing.B) {
	b.Run("n=100", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded100)
		}
	})
}

func BenchmarkEncodeCoords_Allocs(b *testing.B) {
	b.Run("n=100", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords100)
		}
	})
}

// ============================================================================
// Pre-allocated Buffer Benchmarks
// ============================================================================

func BenchmarkEncodeCoords_Preallocated(b *testing.B) {
	b.Run("n=100_no_prealloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords100)
		}
	})

	b.Run("n=100_prealloc", func(b *testing.B) {
		b.ReportAllocs()
		codec := polyline.Codec{Dim: 2, Scale: 1e5}
		buf := make([]byte, 0, 1024)
		for i := 0; i < b.N; i++ {
			_ = codec.EncodeCoords(buf[:0], coords100)
		}
	})
}

// ============================================================================
// Throughput Benchmarks (bytes/sec)
// ============================================================================

func BenchmarkDecodeCoords_Throughput(b *testing.B) {
	b.Run("n=1000", func(b *testing.B) {
		b.SetBytes(int64(len(encoded1000)))
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded1000)
		}
	})

	b.Run("n=10000", func(b *testing.B) {
		b.SetBytes(int64(len(encoded10000)))
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded10000)
		}
	})
}

func BenchmarkEncodeCoords_Throughput(b *testing.B) {
	b.Run("n=1000", func(b *testing.B) {
		result := polyline.EncodeCoords(coords1000)
		b.SetBytes(int64(len(result)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords1000)
		}
	})

	b.Run("n=10000", func(b *testing.B) {
		result := polyline.EncodeCoords(coords10000)
		b.SetBytes(int64(len(result)))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords10000)
		}
	})
}

// ============================================================================
// Streaming Decoder Benchmarks
// ============================================================================

func BenchmarkDecoder2D(b *testing.B) {
	b.Run("n=100", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			dec := polyline.NewDecoder2D(1e5, encoded100)
			for dec.Next() {
				_ = dec.Lat()
				_ = dec.Lon()
			}
		}
	})

	b.Run("n=1000", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			dec := polyline.NewDecoder2D(1e5, encoded1000)
			for dec.Next() {
				_ = dec.Lat()
				_ = dec.Lon()
			}
		}
	})

	b.Run("n=100_reuse", func(b *testing.B) {
		b.ReportAllocs()
		dec := polyline.NewDecoder2D(1e5, nil)
		for i := 0; i < b.N; i++ {
			dec.Reset(encoded100)
			for dec.Next() {
				_ = dec.Lat()
				_ = dec.Lon()
			}
		}
	})
}

// ============================================================================
// Streaming Encoder Benchmarks
// ============================================================================

func BenchmarkEncoder2D(b *testing.B) {
	b.Run("n=100", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			enc := polyline.NewEncoder2D(1e5, 512)
			enc.WriteCoords(coords100)
			_ = enc.Bytes()
		}
	})

	b.Run("n=1000", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			enc := polyline.NewEncoder2D(1e5, 5120)
			enc.WriteCoords(coords1000)
			_ = enc.Bytes()
		}
	})

	b.Run("n=100_reuse", func(b *testing.B) {
		b.ReportAllocs()
		enc := polyline.NewEncoder2D(1e5, 512)
		for i := 0; i < b.N; i++ {
			enc.Reset()
			enc.WriteCoords(coords100)
			_ = enc.Bytes()
		}
	})
}

// ============================================================================
// Comparison: Standard vs Streaming APIs
// ============================================================================

func BenchmarkComparison_Decode_n100(b *testing.B) {
	b.Run("DecodeCoords", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _, _ = polyline.DecodeCoords(encoded100)
		}
	})

	b.Run("Decoder2D", func(b *testing.B) {
		b.ReportAllocs()
		dec := polyline.NewDecoder2D(1e5, nil)
		for i := 0; i < b.N; i++ {
			dec.Reset(encoded100)
			for dec.Next() {
				_ = dec.Lat()
				_ = dec.Lon()
			}
		}
	})

	b.Run("DecodeFlatCoords", func(b *testing.B) {
		b.ReportAllocs()
		codec := polyline.Codec{Dim: 2, Scale: 1e5}
		for i := 0; i < b.N; i++ {
			_, _, _ = codec.DecodeFlatCoords(nil, encoded100)
		}
	})
}

func BenchmarkComparison_Encode_n100(b *testing.B) {
	b.Run("EncodeCoords", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = polyline.EncodeCoords(coords100)
		}
	})

	b.Run("Encoder2D_reuse", func(b *testing.B) {
		b.ReportAllocs()
		enc := polyline.NewEncoder2D(1e5, 512)
		for i := 0; i < b.N; i++ {
			enc.Reset()
			enc.WriteCoords(coords100)
			_ = enc.Bytes()
		}
	})

	b.Run("EncodeFlatCoords", func(b *testing.B) {
		b.ReportAllocs()
		codec := polyline.Codec{Dim: 2, Scale: 1e5}
		for i := 0; i < b.N; i++ {
			_, _ = codec.EncodeFlatCoords(nil, flatCoords100)
		}
	})
}
