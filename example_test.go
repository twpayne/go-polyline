package polyline_test

import (
	"fmt"

	"github.com/twpayne/go-polyline"
)

func ExampleEncodeCoords() {
	coords := [][]float64{
		{38.5, -120.2},
		{40.7, -120.95},
		{43.252, -126.453},
	}
	fmt.Println(string(polyline.EncodeCoords(coords)))
	// Output: _p~iF~ps|U_ulLnnqC_mqNvxq`@
}

func ExampleDecodeCoords() {
	buf := []byte("_p~iF~ps|U_ulLnnqC_mqNvxq`@")
	coords, _, _ := polyline.DecodeCoords(buf)
	fmt.Println(coords)
	// Output: [[38.5 -120.2] [40.7 -120.95] [43.252 -126.453]]
}
