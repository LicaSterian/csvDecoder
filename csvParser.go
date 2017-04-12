package csv

import (
	"fmt"
	"reflect"
	"errors"
	"io"
	"encoding/csv"
	"strconv"
	"time"
)

//Decoder decodes a csv specific io.Reader.
type Decoder struct {
	reader *csv.Reader
	headerKeys map[string]int
}

//NewDecoder receives as an argument a io.Reader and returns a pointer to a CsvDecoder.
func NewDecoder(reader io.Reader) *Decoder {
	r := csv.NewReader(reader)
	return &Decoder{
		reader: r,
	}
}

//Decode method accepts a pointer to a slice with it will populate and returns a error.
//TODO accept tag values as index also
func (p *Decoder) Decode(v interface{}) error {
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
	rte := rt.Elem()
	for {
		line, err := p.reader.Read()
		if err != nil {
			if fmt.Sprint(err) == "EOF" {
				break
			} else {
				return err
			}
		}
		row := reflect.New(rte)
		for i := 0; i < rte.NumField(); i++ {
			var field = rte.Field(i)
			var tag = field.Tag.Get("csv")
			fieldIndex := p.headerKeys[tag]
			var value = line[fieldIndex]
			var settableField = row.Elem().FieldByName(field.Name)
			switch field.Type.Name() {
			case "bool":
				var parsedBool, err = strconv.ParseBool(value)
				if err != nil {
					return err
				}
				settableField.SetBool(parsedBool)
			case "uint", "uint8", "uint16", "uint32", "uint64":
				var parsedUint, err = strconv.ParseUint(value, 10, 64)
				if err != nil {
					return err
				}
				settableField.SetUint(uint64(parsedUint))
			case "int", "int32", "int64":
				var parsedInt, err = strconv.Atoi(value)
				if err != nil {
					return err
				}
				settableField.SetInt(int64(parsedInt))
			case "float32":
				var parsedFloat, err = strconv.ParseFloat(value, 32)
				if err != nil {
					return err
				}
				settableField.SetFloat(parsedFloat)
			case "float64":
				var parsedFloat, err = strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				settableField.SetFloat(parsedFloat)
			case "string":
				settableField.SetString(value)
			case "Time":
				dateTagFormat := field.Tag.Get("csvDate")
				var date, err = time.Parse(dateTagFormat, value)
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