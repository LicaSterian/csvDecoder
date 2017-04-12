package parser

import (
	"testing"
	"bytes"
	"fmt"
)

type CsvRow struct {
	Id        int `csv:"id"`
	FirstName string `csv:"firstName"`
	LastName  string `csv:"lastName"`
}

func TestDecode(t *testing.T) {
	csvContent := `id,firstName,lastName
0,Lica,Sterian
1,Vali,Malinoiu
2,Alex,Leca`

	buf := bytes.NewBufferString(csvContent)
	decoder := NewDecoder(buf)
	var rows []CsvRow
	err := decoder.Decode(&rows)
	if err != nil {
		t.Fatalf("decode error: %s", err.Error())
	}
	fmt.Println("rows:", rows)
}