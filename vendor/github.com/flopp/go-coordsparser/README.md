[![GoDoc](https://godoc.org/github.com/flopp/go-coordsparser?status.svg)](https://godoc.org/github.com/flopp/go-coordsparser)
[![Build Status](https://travis-ci.org/flopp/go-coordsparser.svg)](https://travis-ci.org/flopp/go-coordsparser)
[![Go Report Card](https://goreportcard.com/badge/flopp/go-coordsparser)](https://goreportcard.com/report/flopp/go-coordsparser)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/flopp/go-coordsparser)

# go-coordsparser
A library for parsing (geographic) coordinates in go (golang)

# What?

go-coordsparser allows you to parse lat/lng coordinates from strings in various popular formats. Currently supported formats are:

- **D** (decimal degrees), e.g. `40.76, -73.984`
- **HD** (hemisphere prefix, decimal degrees), e.g. `N 40.76 W 73.984`
- **HDM** (hemisphere prefix, integral degrees, decimal minutes), e.g. `N 40 45.600 W 73 59.040`
- **HDMS** (hemisphere prefix, integral degrees, integral minutes, decimal seconds), e.g. `N 40 45 36.0 W 73 59 02.4`

# How?

### Installing
Installing the library is as easy as

```bash
$ go get github.com/flopp/go-coordsparser
```

The package can then be used through an

```go
import "github.com/flopp/go-coordsparser"
```

### Using
go-coordsparser provides several functions for parsing coordinate strings: a general parsing function `coordsparser.Parse`, which accepts all supported formats, as well as specialized functions `coordsparser.ParseD`, `coordsparser.ParseHD`, `coordsparser.ParseHDM`, `coordsparser.ParseHDMS` for the corresponding coordinate formats.

Each function takes a single string as a parameter and returns an idiomatic `lat, lng, error` triple, where `lat` and `lng` are decimal degrees (`float64`) with -90 ≤ `lat` ≤ 90 and -180 ≤ `lng` ≤ 180.

```go
// parse any format
s1 := "..."
lat1, lng1, err := coordsparser.Parse(s1)
if err != nil {
    fmt.Errorf("Cannot parse coordinates string:", s1)
}

// parse specific format, e.g. HDM
s2 := "..."
lat2, lng2, err = coordsparser.ParseHDM(s2)
if err != nil {
    fmt.Errorf("Cannot parse coordinates string:", s2)
}
```

# License
Copyright 2016 Florian Pigorsch. All rights reserved.

Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
