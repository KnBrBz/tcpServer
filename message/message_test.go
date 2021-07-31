package message

import (
	"strings"
	"testing"
)

func TestMessageContent(t *testing.T) {
	headLength := 2
	cases := map[string]struct {
		Body    []byte
		IsErr   bool
		Content string
	}{
		"empty message":                         {Body: []byte{}, IsErr: true},
		"1 byte message":                        {Body: []byte{0x01}, IsErr: true},
		"2 byte message with zero length":       {Body: []byte{0x00, 0x00}},
		"3 byte message with zero length":       {Body: []byte{0x00, 0x00, 0x02}},
		"empty content is less than length":     {Body: []byte{0x00, 0x01}, IsErr: true},
		"not empty content is less than length": {Body: []byte{0x00, 0x05, 0x41, 0x42, 0x41}, IsErr: true},
		"content is longer than length":         {Body: []byte{0x00, 0x02, 0x41, 0x42, 0x41}, Content: "AB"},
		"content is equal to length":            {Body: append([]byte{0x00, 0xA}, []byte(strings.Repeat("A", 10))...), Content: strings.Repeat("A", 10)},
		"min length":                            {Body: []byte{0x00, 0x01, 0x42}, Content: "B"},
		"max length":                            {Body: append([]byte{0xff, 0xff}, []byte(strings.Repeat("C", 65535))...), Content: strings.Repeat("C", 65535)},
	}
	for title, cs := range cases {
		msg := New(headLength, cs.Body)
		content, err := msg.Content()
		switch {
		case cs.IsErr && err == nil:
			t.Fatalf("case `%s`: error expected", title)
		case !cs.IsErr && err != nil:
			t.Fatalf("case `%s`: expected content `%s`, got error `%v`", title, cs.Content, err)
		case content != cs.Content:
			t.Fatalf("case `%s`: expected content `%s`, got `%s`", title, cs.Content, content)
		}
	}
}
