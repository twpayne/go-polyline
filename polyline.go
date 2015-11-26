// Package polyline implements a Google Maps Encoding Polyline encoder.
//
// See https://developers.google.com/maps/documentation/utilities/polylinealgorithm
package polyline

import (
	"math"
)

func round(x float64) float64 {
	if x < 0 {
		return -math.Floor(-x + 0.5)
	} else {
		return math.Floor(x + 0.5)
	}
}

// A Codec represents an encoder.
type Codec struct {
	Dim   int     // Dimensionality, normally 2
	Scale float64 // Scale, normally 1e5
}

var defaultCodec = Codec{Dim: 2, Scale: 1e5}

// EncodeUint appends the encoding of a single unsigned integer u to result.
func EncodeUint(u uint, result []byte) []byte {
	for u >= 32 {
		result = append(result, byte((u&31)+95))
		u = u >> 5
	}
	result = append(result, byte(u+63))
	return result
}

// EncodeInt appends the encoding of a single signed integer i to result.
func EncodeInt(i int, result []byte) []byte {
	var u uint
	if i < 0 {
		u = uint(^(i << 1))
	} else {
		u = uint(i << 1)
	}
	return EncodeUint(u, result)
}

// EncodeFloat64 appends the encoding of a single float64 f to result.
func (c Codec) EncodeFloat64(f float64, result []byte) []byte {
	return EncodeInt(int(round(c.Scale*f)), result)
}

// EncodeCoords appends the encoding of an array of coordinates coords to result.
func (c Codec) EncodeCoords(coords [][]float64, result []byte) []byte {
	last := make([]float64, c.Dim)
	for _, coord := range coords {
		for i, x := range coord {
			result = c.EncodeFloat64(x-last[i], result)
			last[i] = x
		}
	}
	return result
}

// EncodeFloat64 appends the encoding of a single float64 f to result.
func EncodeFloat64(f float64, result []byte) []byte {
	return defaultCodec.EncodeFloat64(f, result)
}

// EncodeCoords appends the encoding of an array of coordinates coords to result.
func EncodeCoords(coords [][]float64, result []byte) []byte {
	return defaultCodec.EncodeCoords(coords, result)
}
