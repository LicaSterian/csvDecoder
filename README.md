# CSV Decoder [![Build Status](https://travis-ci.org/empatica/csvparser.svg?branch=master)](https://travis-ci.org/empatica/csvparser) [![codecov.io](http://codecov.io/github/empatica/csvparser/coverage.svg?branch=master)](http://codecov.io/github/empatica/csvparser?branch=master)

Simple library that parses a CSV io.Reader and maps it to a given array struct.

### Limitations

- Works only with a struct that contains string, int, uint, bool, time and float fields.

## Getting started

### Install

    go get -u github.com/licasterian/csvdecoder

### Usage

Define your struct:

```go
type YourStruct struct{
  Field1 string
  Field2 int
  Field3 bool
  Field4 float64
  Field5 time.Time `csvDate:"2006-05-07"`
}
```

If you don't add 'csv' tags close to each struct's field, the lib will set the first field using the first column of csv's row, and so on. So the previous struct is the same as:

```go
type Entry struct{
  Field1 string    `csv:"field1"`
  Field2 int       `csv:"field2"`
  Field3 bool      `csv:"field3"`
  Field4 float64   `csv:"field4"`
  Field5 time.Time `csv:"field5" csvDate:"2006-01-02"`
}
```

##### Note for time.Time fields:

It's required to specify a `csvDate` tag that will be used for parsing, following the rules describere [here](http://golang.org/pkg/time/#Parse)

### Parse the file:

```go
decoder := csv.NewDecoder()
var entries []Entry
err := decoder.Decode(&entries)
if err != nil {
    fmt.Printf("decoder.Decode error: %s\n", err.Error())
    return
}
fmt.Printf("entries: %+v\n", entries)
```