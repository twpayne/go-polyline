package polyline

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

func ExampleEncodeCoords() {
	var coords = [][]float64{
		{38.5, -120.2},
		{40.7, -120.95},
		{43.252, -126.453},
	}
	fmt.Printf("%s\n", EncodeCoords(coords))
	// Output: _p~iF~ps|U_ulLnnqC_mqNvxq`@
}

func ExampleDecodeCoords() {
	buf := []byte("_p~iF~ps|U_ulLnnqC_mqNvxq`@")
	coords, _, _ := DecodeCoords(buf)
	fmt.Printf("%v\n", coords)
	// Output: [[38.5 -120.2] [40.7 -120.95] [43.252 -126.453]]
}

func TestUint(t *testing.T) {
	for _, tc := range []struct {
		u uint
		s string
	}{
		{u: 0, s: "?"},
		{u: 31, s: "^"},
		{u: 32, s: "_@"},
		{u: 174, s: "mD"},
	} {
		if got, b, err := DecodeUint([]byte(tc.s)); got != tc.u || len(b) != 0 || err != nil {
			t.Errorf("DecodeUint(%v) = %v, %v, %v, want %v, nil, nil", tc.s, got, err, string(b), tc.u)
		}
		if got := EncodeUint(nil, tc.u); string(got) != tc.s {
			t.Errorf("EncodeUint(%v) = %v, want %v", tc.u, string(got), tc.s)
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
		if _, _, err := DecodeUint([]byte(tc.s)); err == nil || err != tc.err {
			t.Errorf("DecodeUint([]byte(%v)) == _, _, %v, want %v", tc.s, err, tc.err)
		}
		if _, _, err := DecodeInt([]byte(tc.s)); err == nil || err != tc.err {
			t.Errorf("DecodeInt([]byte(%v)) == _, _, %v, want %v", tc.s, err, tc.err)
		}
		if _, _, err := DecodeCoord([]byte(tc.s)); err == nil || err != tc.err {
			t.Errorf("DecodeCoord([]byte(%v)) == _, _, %v, want %v", tc.s, err, tc.err)
		}
		if _, _, err := DecodeCoords([]byte(tc.s)); err == nil || err != tc.err {
			t.Errorf("DecodeCoords([]byte(%v)) == _, _, %v, want %v", tc.s, err, tc.err)
		}
		c := Codec{Dim: 1, Scale: 1e5}
		if _, _, err := c.DecodeFlatCoords([]float64{0}, []byte(tc.s)); err == nil || err != tc.err {
			t.Errorf("DecodeFlatCoords([]byte(%v)) == _, _, %v, want %v", tc.s, err, tc.err)
		}
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
		if _, _, err := DecodeCoords([]byte(tc.s)); err == nil || err != tc.err {
			t.Errorf("DecodeCoords([]byte(%v)) == _, _, %v, want %v", tc.s, err, tc.err)
		}
		c := Codec{Dim: 2, Scale: 1e5}
		if _, _, err := c.DecodeFlatCoords([]float64{0, 0}, []byte(tc.s)); err == nil || err != tc.err {
			t.Errorf("DecodeFlatCoords([]byte(%v)) == _, _, %v, want %v", tc.s, err, tc.err)
		}
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
		if got, b, err := DecodeInt([]byte(tc.s)); got != tc.i || len(b) != 0 || err != nil {
			t.Errorf("DecodeInt(%v) = %v, %v, %v, want %v, nil, nil", tc.s, got, err, string(b), tc.i)
		}
		if got := EncodeInt(nil, tc.i); string(got) != tc.s {
			t.Errorf("EncodeInt(%v) = %v, want %v", tc.i, string(got), tc.s)
		}
	}
}

func TestCoord(t *testing.T) {
	for _, tc := range []struct {
		s string
		c []float64
	}{
		{
			s: "_p~iF~ps|U",
			c: []float64{38.5, -120.2},
		},
	} {
		if got, b, err := DecodeCoord([]byte(tc.s)); !reflect.DeepEqual(got, tc.c) || len(b) != 0 || err != nil {
			t.Errorf("DecodeCoord(%v) = %v, %v, %v, want %v, nil, nil", tc.s, got, err, string(b), tc.c)
		}
		if got := EncodeCoord(tc.c); string(got) != tc.s {
			t.Errorf("EncodeCoord(%v) = %v, want %v", tc.c, string(got), tc.s)
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
		if got, b, err := DecodeCoords([]byte(tc.s)); !reflect.DeepEqual(got, tc.cs) || len(b) != 0 || err != nil {
			t.Errorf("DecodeCoords(%v) = %v, %v, %v, want %v, nil, nil", tc.s, got, string(b), err, tc.cs)
		}
		if got := EncodeCoords(tc.cs); string(got) != tc.s {
			t.Errorf("EncodeCoords(%v) = %v, want %v", tc.cs, string(got), tc.s)
		}
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
		if got, b, err := defaultCodec.DecodeFlatCoords(nil, []byte(tc.s)); !reflect.DeepEqual(got, tc.fcs) || len(b) != 0 || err != nil {
			t.Errorf("defaultCodec.DecodeFlatCoords(nil, %#v) = %v, %v, %v, want %v, nil, nil", tc.s, got, string(b), err, tc.fcs)
		}
		if got, err := defaultCodec.EncodeFlatCoords(nil, tc.fcs); string(got) != tc.s || err != nil {
			t.Errorf("defaultCodec.EncodeFlatCoords(nil, %v) = %v, %v, want %v, nil", tc.fcs, string(got), err, tc.s)
		}
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
		if _, _, err := defaultCodec.DecodeFlatCoords(tc.fcs, []byte(tc.s)); err != tc.err {
			t.Errorf("defaultCodec.DecodeFlatCoords(%v, %v) == _, %v, want _, %v", tc.fcs, tc.s, err, tc.err)
		}
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
		if _, err := defaultCodec.EncodeFlatCoords(nil, tc.fcs); err != tc.err {
			t.Errorf("defaultCodec.EncodeFlatCoords(nil, %v) == _, %v, want _, %v", tc.fcs, err, tc.err)
		}
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
		if got, b, err := tc.c.DecodeCoords([]byte(tc.s)); !reflect.DeepEqual(got, tc.cs) || len(b) != 0 || err != nil {
			t.Errorf("%v.DecodeCoords(%v) = %v, %v, %v, want %v, nil, nil", tc.c, tc.s, got, string(b), err, tc.cs)
		}
		if got := tc.c.EncodeCoords(nil, tc.cs); string(got) != tc.s {
			t.Errorf("%v.EncodeCoords(%v) = %v, want %v", tc.c, tc.cs, string(got), tc.s)
		}
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
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
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
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
