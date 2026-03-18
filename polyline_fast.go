package polyline

// This file contains optimized implementations for the common case of Dim=2.
// These are called from the generic implementations when Dim=2.

// decodeCoordsD2 is an optimized decoder for 2D coordinates.
// It avoids allocations for the coord slice and loop overhead.
func decodeCoordsD2(buf []byte, scale float64) ([][]float64, []byte, error) {
	if len(buf) == 0 {
		return nil, buf, nil
	}

	// Decode first coordinate
	lat, buf, err := DecodeInt(buf)
	if err != nil {
		return nil, nil, err
	}
	lon, buf, err := DecodeInt(buf)
	if err != nil {
		return nil, nil, err
	}

	// Estimate capacity: typical polyline is ~5 bytes per value, 2 values per coord
	estimatedCoords := max(len(buf)/10+1, 8)
	coords := make([][]float64, 0, estimatedCoords)
	coords = append(coords, []float64{float64(lat) / scale, float64(lon) / scale})

	// Decode remaining coordinates (delta-encoded)
	for len(buf) > 0 {
		dLat, remaining, err := DecodeInt(buf)
		if err != nil {
			return nil, nil, err
		}
		dLon, remaining, err := DecodeInt(remaining)
		if err != nil {
			return nil, nil, err
		}
		buf = remaining

		lat += dLat
		lon += dLon
		coords = append(coords, []float64{float64(lat) / scale, float64(lon) / scale})
	}

	return coords, nil, nil
}

// decodeFlatCoordsD2 is an optimized decoder for 2D flat coordinates.
func decodeFlatCoordsD2(flatCoords []float64, buf []byte, scale float64) ([]float64, []byte, error) {
	if len(flatCoords)%2 != 0 {
		return nil, nil, ErrDimensionalMismatch
	}

	var lastLat, lastLon int

	for len(buf) > 0 {
		dLat, remaining, err := DecodeInt(buf)
		if err != nil {
			return nil, nil, err
		}
		dLon, remaining, err := DecodeInt(remaining)
		if err != nil {
			return nil, nil, err
		}
		buf = remaining

		lastLat += dLat
		lastLon += dLon
		flatCoords = append(flatCoords, float64(lastLat)/scale, float64(lastLon)/scale)
	}

	return flatCoords, nil, nil
}

// encodeCoordsD2 is an optimized encoder for 2D coordinates.
func encodeCoordsD2(buf []byte, coords [][]float64, scale float64) []byte {
	if len(coords) == 0 {
		return buf
	}

	// Pre-allocate buffer: estimate ~5 bytes per coordinate value
	if cap(buf)-len(buf) < len(coords)*10 {
		newBuf := make([]byte, len(buf), len(buf)+len(coords)*10)
		copy(newBuf, buf)
		buf = newBuf
	}

	var lastLat, lastLon int

	for _, coord := range coords {
		// Bounds check elimination
		_ = coord[1]

		lat := round(scale * coord[0])
		lon := round(scale * coord[1])

		buf = EncodeInt(buf, lat-lastLat)
		buf = EncodeInt(buf, lon-lastLon)

		lastLat = lat
		lastLon = lon
	}

	return buf
}

// encodeFlatCoordsD2 is an optimized encoder for 2D flat coordinates.
func encodeFlatCoordsD2(buf []byte, flatCoords []float64, scale float64) ([]byte, error) {
	if len(flatCoords)%2 != 0 {
		return nil, ErrDimensionalMismatch
	}

	// Pre-allocate buffer
	if cap(buf)-len(buf) < len(flatCoords)*5 {
		newBuf := make([]byte, len(buf), len(buf)+len(flatCoords)*5)
		copy(newBuf, buf)
		buf = newBuf
	}

	var lastLat, lastLon int

	for i := 0; i < len(flatCoords); i += 2 {
		// Bounds check elimination
		_ = flatCoords[i+1]

		lat := round(scale * flatCoords[i])
		lon := round(scale * flatCoords[i+1])

		buf = EncodeInt(buf, lat-lastLat)
		buf = EncodeInt(buf, lon-lastLon)

		lastLat = lat
		lastLon = lon
	}

	return buf, nil
}
