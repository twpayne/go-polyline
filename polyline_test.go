package polyline

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestUint(t *testing.T) {
	for _, tc := range []struct {
		u            uint
		s            string
		nonCanonical bool
	}{
		{u: 0, s: "?"},
		{u: 2, s: "a?", nonCanonical: true},
		{u: 2, s: "A"},
		{u: 31, s: "^"},
		{u: 32, s: "_@"},
		{u: 174, s: "mD"},
	} {
		got, b, err := DecodeUint([]byte(tc.s))
		assert.NoError(t, err)
		assert.Equal(t, tc.u, got)
		assert.Empty(t, b)
		if !tc.nonCanonical {
			assert.Equal(t, []byte(tc.s), EncodeUint(nil, tc.u))
		}
	}
}

func TestDecodeErrors(t *testing.T) {
	for _, tc := range []struct {
		s   string
		err error
	}{
		{s: ">", err: errInvalidByte},
		{s: "\x80", err: errInvalidByte},
		{s: "_", err: errUnterminatedSequence},
	} {
		var err error
		_, _, err = DecodeUint([]byte(tc.s))
		assert.Equal(t, tc.err, err)
		_, _, err = DecodeInt([]byte(tc.s))
		assert.Equal(t, tc.err, err)
		_, _, err = DecodeCoord([]byte(tc.s))
		assert.Equal(t, tc.err, err)
		_, _, err = DecodeCoords([]byte(tc.s))
		assert.Equal(t, tc.err, err)
		c := Codec{Dim: 1, Scale: 1e5}
		_, _, err = c.DecodeFlatCoords([]float64{0}, []byte(tc.s))
		assert.Equal(t, tc.err, err)
	}
}

func TestMultidimensionalDecodeErrors(t *testing.T) {
	for _, tc := range []struct {
		s   string
		err error
	}{
		{s: "_p~iF~ps|U_p~iF>", err: errInvalidByte},
		{s: "_p~iF~ps|U_p~iF\x80", err: errInvalidByte},
		{s: "_p~iF~ps|U_p~iF~ps|", err: errUnterminatedSequence},
	} {
		_, _, err := DecodeCoords([]byte(tc.s))
		assert.Equal(t, tc.err, err)
		c := Codec{Dim: 2, Scale: 1e5}
		_, _, err = c.DecodeFlatCoords([]float64{0, 0}, []byte(tc.s))
		assert.Equal(t, tc.err, err)
	}
}

func TestInt(t *testing.T) {
	for _, tc := range []struct {
		i int
		s string
	}{
		{i: 3850000, s: "_p~iF"},
		{i: -12020000, s: "~ps|U"},
		{i: -17998321, s: "`~oia@"},
		{i: 220000, s: "_ulL"},
		{i: -75000, s: "nnqC"},
		{i: 255200, s: "_mqN"},
		{i: -550300, s: "vxq`@"},
	} {
		got, b, err := DecodeInt([]byte(tc.s))
		assert.NoError(t, err)
		assert.Empty(t, b)
		assert.Equal(t, tc.i, got)
		assert.Equal(t, []byte(tc.s), EncodeInt(nil, tc.i))
	}
}

func TestCoord(t *testing.T) {
	for _, tc := range []struct {
		s            string
		c            []float64
		nonCanonical bool
	}{
		{
			s: "_p~iF~ps|U",
			c: []float64{38.5, -120.2},
		},
		{
			s:            "a?Z",
			c:            []float64{1e-05, -0.00014},
			nonCanonical: true,
		},
	} {
		got, b, err := DecodeCoord([]byte(tc.s))
		assert.NoError(t, err)
		assert.Empty(t, b)
		assert.Equal(t, tc.c, got)
		if !tc.nonCanonical {
			assert.Equal(t, []byte(tc.s), EncodeCoord(tc.c))
		}
	}
}

func TestCoords(t *testing.T) {
	for _, tc := range []struct {
		cs [][]float64
		s  string
	}{
		{
			cs: [][]float64{{38.5, -120.2}, {40.7, -120.95}, {43.252, -126.453}},
			s:  "_p~iF~ps|U_ulLnnqC_mqNvxq`@",
		},
	} {
		got, b, err := DecodeCoords([]byte(tc.s))
		assert.NoError(t, err)
		assert.Empty(t, b)
		assert.Equal(t, tc.cs, got)
		assert.Equal(t, []byte(tc.s), EncodeCoords(tc.cs))
	}
}

func TestFlatCoords(t *testing.T) {
	for _, tc := range []struct {
		fcs []float64
		s   string
	}{
		{
			fcs: []float64{38.5, -120.2, 40.7, -120.95, 43.252, -126.453},
			s:   "_p~iF~ps|U_ulLnnqC_mqNvxq`@",
		},
	} {
		gotFCS, b, err := defaultCodec.DecodeFlatCoords(nil, []byte(tc.s))
		assert.NoError(t, err)
		assert.Empty(t, b)
		assert.Equal(t, tc.fcs, gotFCS)
		gotBytes, err := defaultCodec.EncodeFlatCoords(nil, tc.fcs)
		assert.NoError(t, err)
		assert.Equal(t, []byte(tc.s), gotBytes)
	}
}

func TestDecodeFlatCoordsErrors(t *testing.T) {
	for _, tc := range []struct {
		fcs []float64
		s   string
		err error
	}{
		{
			fcs: []float64{0},
			s:   "",
			err: errDimensionalMismatch,
		},
		{
			fcs: []float64{0},
			s:   "_p~iF~ps|U",
			err: errDimensionalMismatch,
		},
		{
			fcs: []float64{},
			s:   "_p~iF~ps|U_p~iF",
			err: errUnterminatedSequence,
		},
	} {
		_, _, err := defaultCodec.DecodeFlatCoords(tc.fcs, []byte(tc.s))
		assert.Equal(t, tc.err, err)
	}
}

func TestEncodeFlatCoordErrors(t *testing.T) {
	for _, tc := range []struct {
		fcs []float64
		err error
	}{
		{
			fcs: []float64{0},
			err: errDimensionalMismatch,
		},
	} {
		_, err := defaultCodec.EncodeFlatCoords(nil, tc.fcs)
		assert.Equal(t, tc.err, err)
	}
}

func TestCodec(t *testing.T) {
	for _, tc := range []struct {
		c  Codec
		cs [][]float64
		s  string
	}{
		{
			c:  Codec{Dim: 2, Scale: 1e5},
			cs: [][]float64{{38.5, -120.2}, {40.7, -120.95}, {43.252, -126.453}},
			s:  "_p~iF~ps|U_ulLnnqC_mqNvxq`@",
		},
		{
			c:  Codec{Dim: 2, Scale: 1e6},
			cs: [][]float64{{38.5, -120.2}, {40.7, -120.95}, {43.252, -126.453}},
			s:  "_izlhA~rlgdF_{geC~ywl@_kwzCn`{nI",
		},
	} {
		got, b, err := tc.c.DecodeCoords([]byte(tc.s))
		assert.NoError(t, err)
		assert.Equal(t, tc.cs, got)
		assert.Empty(t, b)
		assert.Equal(t, []byte(tc.s), tc.c.EncodeCoords(nil, tc.cs))
	}
}

func float64ArrayWithin(a, b []float64, prec float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, xa := range a {
		if math.Abs(xa-b[i]) > prec {
			return false
		}
	}
	return true
}

type QuickCoords [][]float64

func (qc QuickCoords) Generate(r *rand.Rand, size int) reflect.Value {
	result := make([][]float64, size)
	for i := range result {
		result[i] = []float64{180*r.Float64() - 90, 360*r.Float64() - 180}
	}
	return reflect.ValueOf(result)
}

func TestCoordsQuick(t *testing.T) {
	f := func(qc QuickCoords) bool {
		buf := EncodeCoords([][]float64(qc))
		cs, buf, err := DecodeCoords(buf)
		if len(buf) != 0 || err != nil {
			return false
		}
		if len(cs) != len(qc) {
			return false
		}
		for i, c := range cs {
			if !float64ArrayWithin(c, qc[i], 5e-6) {
				return false
			}
		}
		return true
	}
	assert.NoError(t, quick.Check(f, nil))
}

type QuickFlatCoords []float64

func (qfc QuickFlatCoords) Generate(r *rand.Rand, size int) reflect.Value {
	result := make([]float64, 2*size)
	for i := range result {
		if i%2 == 0 {
			result[i] = 180*r.Float64() - 90
		} else {
			result[i] = 360*r.Float64() - 180
		}
	}
	return reflect.ValueOf(result)
}

func TestFlatCoordsQuick(t *testing.T) {
	f := func(fqc QuickFlatCoords) bool {
		buf, err := defaultCodec.EncodeFlatCoords(nil, []float64(fqc))
		if err != nil {
			return false
		}
		fcs, _, err := defaultCodec.DecodeFlatCoords(nil, buf)
		if err != nil {
			return false
		}
		return float64ArrayWithin([]float64(fqc), fcs, 5e-6)
	}
	assert.NoError(t, quick.Check(f, nil))
}
