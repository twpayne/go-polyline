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

func EncodeUint(u uint, result []byte) []byte {
	for u > 32 {
		result = append(result, byte((u&31)+95))
		u = u >> 5
	}
	result = append(result, byte(u+63))
	return result
}

func EncodeInt(i int, result []byte) []byte {
	var u uint
	if i < 0 {
		u = uint(^(i << 1))
	} else {
		u = uint(i << 1)
	}
	return EncodeUint(u, result)
}

func EncodeFloat64(f float64, result []byte) []byte {
	return EncodeInt(int(round(1e5*f)), result)
}

func EncodeCoords(coords [][]float64, result []byte) []byte {
	lastLat, lastLong := 0, 0
	for _, coord := range coords {
		lat, long := int(1e5*coord[0]), int(1e5*coord[1])
		result = EncodeInt(lat-lastLat, result)
		result = EncodeInt(long-lastLong, result)
		lastLat, lastLong = lat, long
	}
	return result
}
