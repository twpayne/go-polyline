// Package polyline implements a Google Maps Encoding Polyline encoder and
// decoder. See
// https://developers.google.com/maps/documentation/utilities/polylinealgorithm.
//
// The default codec encodes and decodes two-dimensional coordinates scaled by
// 1e5. For other dimensionalities and scales create a custom Codec.
//
// The package operates on byte slices. Encoding functions take an existing byte
// slice as input (which can be nil) and return a new byte slice with the
// encoded value appended to it, similarly to how Go's append function works. To
// increase performance, you can pre-allocate byte slices, for example by
// passing make([]byte, 0, 128) as the input byte slice. Similarly, decoding
// functions take a byte slice as input and return the remaining unconsumed
// bytes as output.
package polyline

import (
	"errors"
	"math"
	"strconv"
)

// Errors.
var (
	ErrDimensionalMismatch  = errors.New("dimensional mismatch")
	ErrEmpty                = errors.New("empty")
	ErrInvalidByte          = errors.New("invalid byte")
	ErrOverflow             = errors.New("overflow")
	ErrUnterminatedSequence = errors.New("unterminated sequence")
)

func round(x float64) int {
	return int(math.Round(x))
}

// A Codec represents an encoder.
type Codec struct {
	Dim   int     // Dimensionality, normally 2
	Scale float64 // Scale, normally 1e5
}

var defaultCodec = Codec{Dim: 2, Scale: 1e5}

// DecodeUint decodes a single unsigned integer from buf. It returns the decoded
// uint, the remaining unconsumed bytes of buf, and any error.
func DecodeUint(buf []byte) (uint, []byte, error) {
	if len(buf) == 0 {
		return 0, nil, ErrEmpty
	}

	// Fast path: 1-byte value (0-31) - most common for small deltas
	b0 := buf[0]
	if b0 >= 63 && b0 < 95 {
		return uint(b0 - 63), buf[1:], nil
	}
	if b0 < 95 || b0 >= 127 {
		return 0, nil, ErrInvalidByte
	}

	// Fast path: 2-byte value (32-1023)
	if len(buf) >= 2 {
		b1 := buf[1]
		if b1 >= 63 && b1 < 95 {
			return uint(b0-95) | uint(b1-63)<<5, buf[2:], nil
		}
		if b1 < 95 || b1 >= 127 {
			return 0, nil, ErrInvalidByte
		}

		// Fast path: 3-byte value (1024-32767)
		if len(buf) >= 3 {
			b2 := buf[2]
			if b2 >= 63 && b2 < 95 {
				return uint(b0-95) | uint(b1-95)<<5 | uint(b2-63)<<10, buf[3:], nil
			}
			if b2 < 95 || b2 >= 127 {
				return 0, nil, ErrInvalidByte
			}

			// Fast path: 4-byte value
			if len(buf) >= 4 {
				b3 := buf[3]
				if b3 >= 63 && b3 < 95 {
					return uint(b0-95) | uint(b1-95)<<5 | uint(b2-95)<<10 | uint(b3-63)<<15, buf[4:], nil
				}
				if b3 < 95 || b3 >= 127 {
					return 0, nil, ErrInvalidByte
				}
			}
		}
	}

	// Slow path for longer values (5+ bytes)
	return decodeUintSlow(buf)
}

// decodeUintSlow handles decoding of values that require 5+ bytes.
func decodeUintSlow(buf []byte) (uint, []byte, error) {
	n := min(strconv.IntSize/5, len(buf))
	// Bounds check elimination hint
	_ = buf[n-1]

	var u, shift uint
	for i := range n {
		b := buf[i]
		if b >= 95 && b < 127 {
			u += (uint(b) - 95) << shift
			shift += 5
		} else if b >= 63 && b < 95 {
			u += (uint(b) - 63) << shift
			return u, buf[i+1:], nil
		} else {
			return 0, nil, ErrInvalidByte
		}
	}
	if len(buf) <= strconv.IntSize/5 {
		return 0, nil, ErrUnterminatedSequence
	}
	maxDigit := byte(1<<(strconv.IntSize-5*(strconv.IntSize/5)) - 1)
	b := buf[n]
	if b >= 63 && b <= 63+maxDigit {
		u += (uint(b) - 63) << shift
		return u, buf[n+1:], nil
	}
	if b < 127 {
		return 0, nil, ErrOverflow
	}
	return 0, nil, ErrInvalidByte
}

// DecodeInt decodes a single signed integer from buf. It returns the decoded
// int, the remaining unconsumed bytes of buf, and any error.
func DecodeInt(buf []byte) (int, []byte, error) {
	switch u, buf, err := DecodeUint(buf); {
	case err != nil:
		return 0, nil, err
	case u&1 == 0:
		return int(u >> 1), buf, nil
	case u == math.MaxUint64:
		return math.MinInt64, buf, nil
	default:
		return -int((u + 1) >> 1), buf, nil
	}
}

// EncodeUint appends the encoding of a single unsigned integer u to buf and
// returns the new buf.
func EncodeUint(buf []byte, u uint) []byte {
	// Fast path for small values (most common case: deltas between coordinates)
	if u < 32 {
		return append(buf, byte(u+63))
	}
	if u < 1024 {
		return append(buf, byte((u&31)+95), byte((u>>5)+63))
	}
	if u < 32768 {
		return append(buf, byte((u&31)+95), byte(((u>>5)&31)+95), byte((u>>10)+63))
	}
	if u < 1048576 {
		return append(buf, byte((u&31)+95), byte(((u>>5)&31)+95), byte(((u>>10)&31)+95), byte((u>>15)+63))
	}
	// Slow path for larger values
	for u >= 32 {
		buf = append(buf, byte((u&31)+95))
		u >>= 5
	}
	return append(buf, byte(u+63))
}

// EncodeInt appends the encoding of a single signed integer i to buf and
// returns the new buf.
func EncodeInt(buf []byte, i int) []byte {
	var u uint
	if i < 0 {
		u = uint(^(i << 1))
	} else {
		u = uint(i << 1)
	}
	return EncodeUint(buf, u)
}

// DecodeCoord decodes a single coordinate from buf. It returns the coordinate,
// the remaining unconsumed bytes of buf, and any error.
func (c Codec) DecodeCoord(buf []byte) ([]float64, []byte, error) {
	coord := make([]float64, c.Dim)
	for i := range coord {
		var err error
		var j int
		j, buf, err = DecodeInt(buf)
		if err != nil {
			return nil, nil, err
		}
		coord[i] = float64(j) / c.Scale
	}
	return coord, buf, nil
}

// DecodeCoords decodes an array of coordinates from buf. It returns the
// coordinates, the remaining unconsumed bytes of buf, and any error.
func (c Codec) DecodeCoords(buf []byte) ([][]float64, []byte, error) {
	// Fast path for the common 2D case
	if c.Dim == 2 {
		return decodeCoordsD2(buf, c.Scale)
	}

	if len(buf) == 0 {
		return nil, buf, nil
	}
	var coord []float64
	var err error
	coord, buf, err = c.DecodeCoord(buf)
	if err != nil {
		return nil, nil, err
	}
	coords := [][]float64{coord}
	for i := 1; len(buf) > 0; i++ {
		coord, buf, err = c.DecodeCoord(buf)
		if err != nil {
			return nil, nil, err
		}
		for j := range coord {
			coord[j] += coords[i-1][j]
		}
		coords = append(coords, coord)
	}
	return coords, nil, nil
}

// DecodeFlatCoords decodes coordinates from buf, appending them to a
// one-dimensional array. It returns the coordinates, the remaining unconsumed
// bytes in buf, and any error.
func (c Codec) DecodeFlatCoords(flatCoords []float64, buf []byte) ([]float64, []byte, error) {
	// Fast path for the common 2D case
	if c.Dim == 2 {
		return decodeFlatCoordsD2(flatCoords, buf, c.Scale)
	}

	if len(flatCoords)%c.Dim != 0 {
		return nil, nil, ErrDimensionalMismatch
	}
	last := make([]int, c.Dim)
	for len(buf) > 0 {
		for j := range c.Dim {
			var err error
			var k int
			k, buf, err = DecodeInt(buf)
			if err != nil {
				return nil, nil, err
			}
			last[j] += k
			flatCoords = append(flatCoords, float64(last[j])/c.Scale)
		}
	}
	return flatCoords, nil, nil
}

// EncodeCoord encodes a single coordinate to buf and returns the new buf.
func (c Codec) EncodeCoord(buf []byte, coord []float64) []byte {
	for _, x := range coord {
		buf = EncodeInt(buf, round(c.Scale*x))
	}
	return buf
}

// EncodeCoords appends the encoding of an array of coordinates coords to buf
// and returns the new buf.
func (c Codec) EncodeCoords(buf []byte, coords [][]float64) []byte {
	// Fast path for the common 2D case
	if c.Dim == 2 {
		return encodeCoordsD2(buf, coords, c.Scale)
	}

	last := make([]int, c.Dim)
	for _, coord := range coords {
		for i, x := range coord {
			ex := round(c.Scale * x)
			buf = EncodeInt(buf, ex-last[i])
			last[i] = ex
		}
	}
	return buf
}

// EncodeFlatCoords encodes a one-dimensional array of coordinates to buf. It
// returns the new buf and any error.
func (c Codec) EncodeFlatCoords(buf []byte, flatCoords []float64) ([]byte, error) {
	// Fast path for the common 2D case
	if c.Dim == 2 {
		return encodeFlatCoordsD2(buf, flatCoords, c.Scale)
	}

	if len(flatCoords)%c.Dim != 0 {
		return nil, ErrDimensionalMismatch
	}
	last := make([]int, c.Dim)
	for i, x := range flatCoords {
		ex := round(c.Scale * x)
		j := i % c.Dim
		buf = EncodeInt(buf, ex-last[j])
		last[j] = ex
	}
	return buf, nil
}

// DecodeCoord decodes a single coordinate from buf using the default codec. It
// returns the coordinate, the remaining bytes in buf, and any error.
func DecodeCoord(buf []byte) ([]float64, []byte, error) {
	return defaultCodec.DecodeCoord(buf)
}

// DecodeCoords decodes an array of coordinates from buf using the default
// codec. It returns the coordinates, the remaining bytes in buf, and any error.
func DecodeCoords(buf []byte) ([][]float64, []byte, error) {
	return defaultCodec.DecodeCoords(buf)
}

// EncodeCoord returns the encoding of an array of coordinates using the default
// codec.
func EncodeCoord(coord []float64) []byte {
	return defaultCodec.EncodeCoord(nil, coord)
}

// EncodeCoords returns the encoding of an array of coordinates using the
// default codec.
func EncodeCoords(coords [][]float64) []byte {
	return defaultCodec.EncodeCoords(nil, coords)
}
