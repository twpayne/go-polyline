package polyline

// This file provides streaming encoder and decoder APIs that minimize allocations
// by reusing internal buffers and allowing callers to provide their own storage.

// Decoder is a streaming decoder for polyline-encoded data.
// It allows decoding coordinates one at a time without allocating
// a new slice for each coordinate.
type Decoder struct {
	c     Codec
	buf   []byte
	last  []int     // Accumulator for delta-decoded values
	coord []float64 // Reusable coordinate buffer
	err   error
}

// NewDecoder creates a new streaming decoder for the given codec and input buffer.
func NewDecoder(c Codec, buf []byte) *Decoder {
	return &Decoder{
		c:     c,
		buf:   buf,
		last:  make([]int, c.Dim),
		coord: make([]float64, c.Dim),
	}
}

// Next advances to the next coordinate. Returns true if a coordinate was decoded,
// false if the input is exhausted or an error occurred.
func (d *Decoder) Next() bool {
	if d.err != nil || len(d.buf) == 0 {
		return false
	}

	for i := range d.c.Dim {
		var delta int
		delta, d.buf, d.err = DecodeInt(d.buf)
		if d.err != nil {
			return false
		}
		d.last[i] += delta
		d.coord[i] = float64(d.last[i]) / d.c.Scale
	}
	return true
}

// Coord returns the current coordinate. The returned slice is reused
// between calls to Next(), so callers must copy it if they need to retain it.
func (d *Decoder) Coord() []float64 {
	return d.coord
}

// Err returns any error that occurred during decoding.
func (d *Decoder) Err() error {
	return d.err
}

// Remaining returns the unconsumed portion of the input buffer.
func (d *Decoder) Remaining() []byte {
	return d.buf
}

// Reset resets the decoder to decode from a new buffer.
func (d *Decoder) Reset(buf []byte) {
	d.buf = buf
	d.err = nil
	for i := range d.last {
		d.last[i] = 0
	}
}

// Decoder2D is an optimized streaming decoder for 2D coordinates.
// It avoids the overhead of the generic Decoder for the common case.
type Decoder2D struct {
	scale    float64
	buf      []byte
	lastLat  int
	lastLon  int
	lat, lon float64
	err      error
}

// NewDecoder2D creates a new streaming decoder for 2D coordinates.
func NewDecoder2D(scale float64, buf []byte) *Decoder2D {
	return &Decoder2D{
		scale: scale,
		buf:   buf,
	}
}

// Next advances to the next coordinate. Returns true if a coordinate was decoded.
func (d *Decoder2D) Next() bool {
	if d.err != nil || len(d.buf) == 0 {
		return false
	}

	var dLat, dLon int
	dLat, d.buf, d.err = DecodeInt(d.buf)
	if d.err != nil {
		return false
	}
	dLon, d.buf, d.err = DecodeInt(d.buf)
	if d.err != nil {
		return false
	}

	d.lastLat += dLat
	d.lastLon += dLon
	d.lat = float64(d.lastLat) / d.scale
	d.lon = float64(d.lastLon) / d.scale
	return true
}

// Lat returns the latitude of the current coordinate.
func (d *Decoder2D) Lat() float64 {
	return d.lat
}

// Lon returns the longitude of the current coordinate.
func (d *Decoder2D) Lon() float64 {
	return d.lon
}

// Coord returns the current coordinate as [lat, lon].
// Note: This allocates a new slice. For zero-allocation access, use Lat() and Lon().
func (d *Decoder2D) Coord() []float64 {
	return []float64{d.lat, d.lon}
}

// Err returns any error that occurred during decoding.
func (d *Decoder2D) Err() error {
	return d.err
}

// Remaining returns the unconsumed portion of the input buffer.
func (d *Decoder2D) Remaining() []byte {
	return d.buf
}

// Reset resets the decoder to decode from a new buffer.
func (d *Decoder2D) Reset(buf []byte) {
	d.buf = buf
	d.err = nil
	d.lastLat = 0
	d.lastLon = 0
}

// Encoder is a streaming encoder for polyline data.
// It allows encoding coordinates one at a time to a reusable buffer.
type Encoder struct {
	c    Codec
	buf  []byte
	last []int
}

// NewEncoder creates a new streaming encoder with the specified initial capacity.
func NewEncoder(c Codec, initialCap int) *Encoder {
	return &Encoder{
		c:    c,
		buf:  make([]byte, 0, initialCap),
		last: make([]int, c.Dim),
	}
}

// WriteCoord encodes a single coordinate and appends it to the internal buffer.
func (e *Encoder) WriteCoord(coord []float64) {
	for i, x := range coord {
		ex := round(e.c.Scale * x)
		e.buf = EncodeInt(e.buf, ex-e.last[i])
		e.last[i] = ex
	}
}

// WriteCoords encodes multiple coordinates and appends them to the internal buffer.
func (e *Encoder) WriteCoords(coords [][]float64) {
	for _, coord := range coords {
		e.WriteCoord(coord)
	}
}

// Bytes returns the encoded polyline. The returned slice is owned by the Encoder
// and will be overwritten by subsequent calls to Reset().
func (e *Encoder) Bytes() []byte {
	return e.buf
}

// Len returns the current length of the encoded data.
func (e *Encoder) Len() int {
	return len(e.buf)
}

// Reset clears the encoder's buffer and state for reuse.
func (e *Encoder) Reset() {
	e.buf = e.buf[:0]
	for i := range e.last {
		e.last[i] = 0
	}
}

// Encoder2D is an optimized streaming encoder for 2D coordinates.
type Encoder2D struct {
	scale            float64
	buf              []byte
	lastLat, lastLon int
}

// NewEncoder2D creates a new streaming encoder for 2D coordinates.
func NewEncoder2D(scale float64, initialCap int) *Encoder2D {
	return &Encoder2D{
		scale: scale,
		buf:   make([]byte, 0, initialCap),
	}
}

// WriteCoord encodes a lat/lon coordinate.
func (e *Encoder2D) WriteCoord(lat, lon float64) {
	latInt := round(e.scale * lat)
	lonInt := round(e.scale * lon)
	e.buf = EncodeInt(e.buf, latInt-e.lastLat)
	e.buf = EncodeInt(e.buf, lonInt-e.lastLon)
	e.lastLat = latInt
	e.lastLon = lonInt
}

// WriteCoordSlice encodes a coordinate from a [lat, lon] slice.
func (e *Encoder2D) WriteCoordSlice(coord []float64) {
	_ = coord[1] // bounds check
	e.WriteCoord(coord[0], coord[1])
}

// WriteCoords encodes multiple coordinates.
func (e *Encoder2D) WriteCoords(coords [][]float64) {
	for _, coord := range coords {
		e.WriteCoordSlice(coord)
	}
}

// Bytes returns the encoded polyline.
func (e *Encoder2D) Bytes() []byte {
	return e.buf
}

// Len returns the current length of the encoded data.
func (e *Encoder2D) Len() int {
	return len(e.buf)
}

// Reset clears the encoder for reuse.
func (e *Encoder2D) Reset() {
	e.buf = e.buf[:0]
	e.lastLat = 0
	e.lastLon = 0
}

// DecodeCoordInto decodes a single coordinate into the provided slice.
// The slice must have length >= c.Dim.
func (c Codec) DecodeCoordInto(coord []float64, buf []byte) ([]byte, error) {
	if len(coord) < c.Dim {
		return nil, ErrDimensionalMismatch
	}
	for i := range c.Dim {
		var err error
		var j int
		j, buf, err = DecodeInt(buf)
		if err != nil {
			return nil, err
		}
		coord[i] = float64(j) / c.Scale
	}
	return buf, nil
}
