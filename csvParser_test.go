package csv

import (
	"testing"
	"bytes"
	"time"
)

type Entry struct {
	Id        int `csv:"id"`
	Nickname  string `csv:"nickname"`
	CreatedOn time.Time `csv:"createdOn" csvDate:"2006-01-02"`
}

func TestDecode(t *testing.T) {
	content := `id,nickname,createdOn
0,SLIde,2017-04-12`
	buf := bytes.NewBufferString(content)
	decoder := NewDecoder(buf)
	var rows []Entry
	err := decoder.Decode(&rows)
	if err != nil {
		t.Fatalf("decode error: %s", err.Error())
	}
	if len(rows) != 1 {
		t.Fatalf("number of decoded rows should be 1 but it's: %d", len(rows))
	}
	firstEntry := rows[0]
	if firstEntry.Id != 0 || firstEntry.Nickname != "SLIde" || firstEntry.CreatedOn.Format("2006-01-02") != "2017-04-12" {
		t.Fatalf(`First csv entry should be {Id:0, Nickname:SLIde, CreatedOn: 2017-04-12 00:00:00 +0000UTC}, but it's %+v`, firstEntry)
	}
}