# gopolyline2

Package polyline implements a Google Maps Encoding Polyline encoder.

See https://godoc.org/github.com/twpayne/gopolyline2/polyline.

Example:

```go
func ExampleEncodeCoords() {
	var coords = [][]float64{
		[]float64{38.5, -120.2},
		[]float64{40.7, -120.95},
		[]float64{43.252, -126.453},
	}
	fmt.Printf("%s\n", EncodeCoords(coords, nil))
	// Output: _p~iF~ps|U_ulLnnqC_mqNvxq`@
}
```

[License](LICENSE)
