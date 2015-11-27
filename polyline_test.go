package polyline

import (
	"fmt"
	"reflect"
	"testing"
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
