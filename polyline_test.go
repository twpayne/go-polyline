package polyline

import (
	"fmt"
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

func TestEncodeInt(t *testing.T) {
	for _, c := range []struct {
		i    int
		want string
	}{
		{i: 3850000, want: "_p~iF"},
		{i: -12020000, want: "~ps|U"},
		{i: -17998321, want: "`~oia@"},
	} {
		if got := EncodeInt(c.i, nil); string(got) != c.want {
			t.Errorf("EncodeInt(%v) = %v, want %v", c.i, string(got), c.want)
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
