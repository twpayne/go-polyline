# go-polyline

[![GoDoc](https://godoc.org/github.com/twpayne/go-polyline?status.svg)](https://godoc.org/github.com/twpayne/go-polyline)
[![Coverage Status](https://coveralls.io/repos/github/twpayne/go-polyline/badge.svg)](https://coveralls.io/github/twpayne/go-polyline)

Package `polyline` implements a Google Maps Encoding Polyline encoder and decoder.

## Encoding example

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

## Decoding example

```go
func ExampleDecodeCoords() {
    buf := []byte("_p~iF~ps|U_ulLnnqC_mqNvxq`@")
    coords, _, _ := DecodeCoords(buf)
    fmt.Printf("%v\n", coords)
    // Output: [[38.5 -120.2] [40.7 -120.95] [43.252 -126.453]]
}
```

## License

BSD-2-Clause
