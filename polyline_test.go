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
		[]float64{38.5, -120.2},
		[]float64{40.7, -120.95},
		[]float64{43.252, -126.453},
	}
	fmt.Printf("%s\n", EncodeCoords(coords, nil))
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
		if got := EncodeUint(tc.u, nil); string(got) != tc.s {
			t.Errorf("EncodeUint(%v) = %v, want %v", tc.u, string(got), tc.s)
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
		if got := EncodeInt(tc.i, nil); string(got) != tc.s {
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
		if got := EncodeCoord(tc.c, nil); string(got) != tc.s {
			t.Errorf("EncodeCoord(%v, nil) = %v, want %v", tc.c, got, string(got), tc.s)
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
		if got := EncodeCoords(tc.cs, nil); string(got) != tc.s {
			t.Errorf("EncodeCoords(%v) = %v, want %v", tc.cs, string(got), tc.s)
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
		if got := tc.c.EncodeCoords(tc.cs, nil); string(got) != tc.s {
			t.Errorf("%v.EncodeCoords(%v) = %v, want %v", tc.c, tc.cs, string(got), tc.s)
		}
	}
}

type QuickCoords [][]float64

func (qc QuickCoords) Generate(rand *rand.Rand, size int) reflect.Value {
	result := make([][]float64, size)
	for i := range result {
		result[i] = []float64{180*rand.Float64() - 90, 360*rand.Float64() - 180}
	}
	return reflect.ValueOf(result)
}

func TestCoordsQuick(t *testing.T) {
	f := func(qc QuickCoords) bool {
		buf := EncodeCoords([][]float64(qc), nil)
		cs, buf, err := DecodeCoords(buf)
		if len(buf) != 0 || err != nil {
			return false
		}
		if len(cs) != len(qc) {
			return false
		}
		for i, c := range cs {
			if len(c) != len(qc[i]) {
				return false
			}
			for j, x := range c {
				if math.Abs(x-qc[i][j]) > 5e-6 {
					return false
				}
			}
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
