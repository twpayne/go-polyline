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

func TestUint(t *testing.T) {
	for _, c := range []struct {
		u uint
		s string
	}{
		{u: 0, s: "?"},
		{u: 31, s: "^"},
		{u: 32, s: "_@"},
		{u: 174, s: "mD"},
	} {
		if got, b, err := DecodeUint([]byte(c.s)); got != c.u || len(b) != 0 || err != nil {
			t.Errorf("DecodeUint(%v) = %v, %v, %v, want %v, nil, nil", c.s, got, err, string(b), c.u)
		}
		if got := EncodeUint(c.u, nil); string(got) != c.s {
			t.Errorf("EncodeUint(%v) = %v, want %v", c.u, string(got), c.s)
		}
	}
}

func TestInt(t *testing.T) {
	for _, c := range []struct {
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
		if got, b, err := DecodeInt([]byte(c.s)); got != c.i || len(b) != 0 || err != nil {
			t.Errorf("DecodeInt(%v) = %v, %v, %v, want %v, nil, nil", c.s, got, err, string(b), c.i)
		}
		if got := EncodeInt(c.i, nil); string(got) != c.s {
			t.Errorf("EncodeInt(%v) = %v, want %v", c.i, string(got), c.s)
		}
	}
}

func TestCoord(t *testing.T) {
	for _, c := range []struct {
		s string
		c []float64
	}{
		{
			s: "_p~iF~ps|U",
			c: []float64{38.5, -120.2},
		},
	} {
		if got, b, err := DecodeCoord([]byte(c.s)); !reflect.DeepEqual(got, c.c) || len(b) != 0 || err != nil {
			t.Errorf("DecodeCoord(%v) = %v, %v, %v, want %v, nil, nil", c.s, got, err, string(b), c.c)
		}
		if got := EncodeCoord(c.c, nil); string(got) != c.s {
			t.Errorf("EncodeCoord(%v, nil) = %v, want %v", c.c, got, string(got), c.s)
		}
	}
}

func TestEncodeCoords(t *testing.T) {
	for _, c := range []struct {
		coords [][]float64
		want   string
	}{
		{
			coords: [][]float64{{38.5, -120.2}, {40.7, -120.95}, {43.252, -126.453}},
			want:   "_p~iF~ps|U_ulLnnqC_mqNvxq`@",
		},
	} {
		if got := EncodeCoords(c.coords, nil); string(got) != c.want {
			t.Errorf("EncodeCoords(%v) = %v, want %v", c.coords, string(got), c.want)
		}
	}
}
