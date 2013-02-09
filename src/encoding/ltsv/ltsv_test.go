package ltsv

import (
	"bytes"
	"testing"
	"io"
)

type readerTest struct {
	value  string
	records []map[string]string
}

var readerTests = []readerTest {
	{
		`host:127.0.0.1	ident:-	user:frank	time:[10/Oct/2000:13:55:36 -0700]	req:GET /apache_pb.gif

HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)
`,
		[]map[string]string{
			{"host": "127.0.0.1", "ident": "-", "user": "frank", "time": "[10/Oct/2000:13:55:36 -0700]", "req": "GET /apache_pb.gif"},
			{"status": "200", "size": "2326", "referer": "http://www.example.com/start.html", "ua": "Mozilla/4.08 [en] (Win98; I ;Nav)"},
		},
	},
	{
		` trimspace :こんにちは
		 trim space :こんばんは
日本語:ラベル
nolabelnofield
ha,s.p-un_ct: おはよう `,
		[]map[string]string{
			{"trimspace": "こんにちは"},
			{"trim space": "こんばんは"},
			{"日本語": "ラベル"},
			{"ha,s.p-un_ct": " おはよう "},
		},
	},
}

func TestReader(t *testing.T) {
	for n, test := range readerTests {
		reader := NewReader(bytes.NewBufferString(test.value))
		for i, result := range test.records {
			record, err := reader.Read()
			if err != nil {
				t.Errorf("error %v at test %d, line %d", err, n, i)
			}
			for label, field := range result {
				if record[label] != field {
					t.Errorf("wrong field %s: test %d, line %d, label %s, field %s", record[label], n, i, label, field)
				}
			}
			if len(result) != len(record) {
				t.Errorf("wrong size %d, %v :test %d, line %d", len(record), record, n, i)
			}
		}
		_, err := reader.Read()
		if err != io.EOF {
			t.Errorf("expected EOF: %v", err)
		}
	}
}

func TestWriter(t *testing.T) {
}
