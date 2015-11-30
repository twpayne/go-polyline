# polyline

Package polyline implements a Google Maps Encoding Polyline encoder and decoder.

See https://godoc.org/github.com/twpayne/go-polyline.

Encoding example:

```go
func ExampleEncodeCoords() {
	var coords = [][]float64{
		[]float64{38.5, -120.2},
		[]float64{40.7, -120.95},
		[]float64{43.252, -126.453},
	}
	fmt.Printf("%s\n", EncodeCoords(coords))
	// Output: _p~iF~ps|U_ulLnnqC_mqNvxq`@
}
```

Decoding example:

```go
func ExampleDecodeCoords() {
	buf := []byte("_p~iF~ps|U_ulLnnqC_mqNvxq`@")
	coords, _, _ := DecodeCoords(buf)
	fmt.Printf("%v\n", coords)
	// Output: [[38.5 -120.2] [40.7 -120.95] [43.252 -126.453]]
}
```

[License](LICENSE)
