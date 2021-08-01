package message

import (
	"regexp"
	"strings"
	"testing"
)

func TestMessageContent(t *testing.T) {
	headLength := HeadLength
	tagReg := regexp.MustCompile("#.+#")
	cases := map[string]struct {
		Body  []byte
		IsErr bool
		Tag   string
	}{
		"empty message":                         {Body: []byte{}, IsErr: true},
		"1 byte message":                        {Body: []byte{0x01}, IsErr: true},
		"2 byte message with zero length":       {Body: []byte{0x00, 0x00}},
		"3 byte message with zero length":       {Body: []byte{0x00, 0x00, 0x02}},
		"empty content is less than length":     {Body: []byte{0x00, 0x01}, IsErr: true},
		"not empty content is less than length": {Body: []byte{0x00, 0x05, 0x41, 0x42, 0x41}, IsErr: true},
		"content is longer than length":         {Body: []byte{0x00, 0x02, 0x41, 0x42, 0x41}, IsErr: true},
		"content is equal to length":            {Body: append([]byte{0x00, 0xA}, []byte(strings.Repeat("A", 10))...)},
		"min length":                            {Body: []byte{0x00, 0x01, 0x42}},
		"max length":                            {Body: append([]byte{0xff, 0xff}, []byte(strings.Repeat("C", 65535))...)},
		"tag at start":                          {Body: append([]byte{0x00, 0xA}, []byte("#test#some")...), Tag: "#test#"},
		"tag in the middle":                     {Body: append([]byte{0x00, 0xA}, []byte("so#test#me")...)},
		"tag in the end":                        {Body: append([]byte{0x00, 0xA}, []byte("some#test#")...)},
		"empty tag":                             {Body: append([]byte{0x00, 0x06}, []byte("##some")...)},
	}
	for title, cs := range cases {
		msg := New(headLength, cs.Body)
		err := msg.Validate(tagReg)
		tag := msg.Tag()
		switch {
		case cs.IsErr && err == nil:
			t.Fatalf("case `%s`: error expected", title)
		case !cs.IsErr && err != nil:
			t.Fatalf("case `%s`: expected tag `%s`, got error `%v`", title, cs.Tag, err)
		case tag != cs.Tag:
			t.Fatalf("case `%s`: expected tag `%s`, got `%s`", title, cs.Tag, tag)
		}
	}
}
