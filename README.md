# polyline

[![Build Status](https://travis-ci.org/twpayne/go-polyline.svg?branch=master)](https://travis-ci.org/twpayne/go-polyline)
[![GoDoc](https://godoc.org/github.com/twpayne/go-polyline?status.svg)](https://godoc.org/github.com/twpayne/go-polyline)
[![Report Card](https://goreportcard.com/badge/github.com/twpayne/go-polyline)](https://goreportcard.com/report/github.com/twpayne/go-polyline)

Package polyline implements a Google Maps Encoding Polyline encoder and decoder.

Encoding example:

```go
func ExampleEncodeCoords() {
	var coords = [][]float64{
		{38.5, -120.2},
		{40.7, -120.95},
		{43.252, -126.453},
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
