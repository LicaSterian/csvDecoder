package parser

import (
	"fmt"
	"reflect"
	"errors"
	"io"
	"encoding/csv"
	"strconv"
	"time"
)

//Parser parse a csv file and returns an array of pointers of the type specified
type Decoder interface {
	Decode(v interface{}) error
}

//CsvDecoder parses a csv file and returns an array of pointers the type specified
type CsvDecoder struct {
	reader *csv.Reader
	//SkipEmptyValues bool
	headerKeys map[string]int
}

func NewDecoder(reader io.Reader) *CsvDecoder {
	csvReader := csv.NewReader(reader)
	return &CsvDecoder{
		reader: csvReader,
	}
}

//Parse creates the array of the given type from the csv file
func (p *CsvDecoder) Decode(v interface{}) error {
	rv := reflect.ValueOf(v)
	re := rv.Elem()
	rt := rv.Type()
	rk := rt.Kind()
	if rk == reflect.Ptr && !rv.IsNil() {
		rv = reflect.Indirect(rv)
		rt = rv.Type()
		rk = rt.Kind()
	} else {
		return errors.New("Decode parameter must be a non-nil pointer to a slice")
	}
	headerRecords, err := p.reader.Read()
	if err != nil {
		return err
	}
	p.headerKeys = map[string]int{}
	for headerRecordIndex, headerRecord := range headerRecords {
		p.headerKeys[headerRecord] = headerRecordIndex
	}
	resultElem := rt.Elem()
	for {
		line, err := p.reader.Read()
		if err != nil {
			if fmt.Sprint(err) == "EOF" {
				break
			} else {
				return err
			}
		}
		row := reflect.New(resultElem)
		for fieldIndex := 0; fieldIndex < resultElem.NumField(); fieldIndex++ {
			var currentField = resultElem.Field(fieldIndex)
			var csvTag = currentField.Tag.Get("csv")
			csvFieldIndex := p.headerKeys[csvTag]
			var csvValue = line[csvFieldIndex]
			var settableField = row.Elem().FieldByName(currentField.Name)
			switch currentField.Type.Name() {
			case "bool":
				var parsedBool, err = strconv.ParseBool(csvValue)
				if err != nil {
					return err
				}
				settableField.SetBool(parsedBool)
			case "uint", "uint8", "uint16", "uint32", "uint64":
				var parsedUint, err = strconv.ParseUint(csvValue, 10, 64)
				if err != nil {
					return err
				}
				settableField.SetUint(uint64(parsedUint))
			case "int", "int32", "int64":
				var parsedInt, err = strconv.Atoi(csvValue)
				if err != nil {
					return err
				}
				settableField.SetInt(int64(parsedInt))
			case "float32":
				var parsedFloat, err = strconv.ParseFloat(csvValue, 32)
				if err != nil {
					return err
				}
				settableField.SetFloat(parsedFloat)
			case "float64":
				var parsedFloat, err = strconv.ParseFloat(csvValue, 64)
				if err != nil {
					return err
				}
				settableField.SetFloat(parsedFloat)
			case "string":
				settableField.SetString(csvValue)
			case "Time":
				var date, err = time.Parse(currentField.Tag.Get("csvDate"), csvValue)
				if err != nil {
					return err
				}
				settableField.Set(reflect.ValueOf(date))
			}
		}
		re.Set(reflect.Append(re, row.Elem()))
	}
	return nil
}
