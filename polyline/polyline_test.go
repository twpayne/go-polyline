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

func TestEncodeUint(t *testing.T) {
	for _, c := range []struct {
		u    uint
		want string
	}{
		{u: 0, want: "?"},
		{u: 31, want: "^"},
		{u: 32, want: "_@"},
		{u: 174, want: "mD"},
	} {
		if got := EncodeUint(c.u, nil); string(got) != c.want {
			t.Errorf("EncodeUint(%v) = %v, want %v", c.u, string(got), c.want)
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
