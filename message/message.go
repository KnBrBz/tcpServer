package message

import (
	"encoding/binary"
	"regexp"

	"github.com/pkg/errors"
)

type M struct {
	body       []byte
	headLength int
	tag        string
	isValid    bool
}

func New(headLength int, body []byte) *M {
	dest := make([]byte, len(body))
	copy(dest, body)
	return &M{
		body:       dest,
		headLength: headLength,
	}
}

func (m *M) Body() []byte {
	return m.body
}

func (m *M) Tag() string {
	return m.tag
}

func (m *M) Validate(reg *regexp.Regexp) (err error) {
	const funcTitle = packageTitle + "*M.validate"
	if m.isValid {
		return
	}
	contentLength, err := m.length()
	if err != nil {
		return errors.Wrapf(err, "%s: wrong message format", funcTitle)
	}
	if contentLength == 0 {
		return
	}
	totalLength := m.headLength + contentLength
	if len(m.body) != totalLength {
		return errors.Wrapf(errors.New("wrong content length"), "%s: wrong message format", funcTitle)
	}
	m.isValid = true
	if reg != nil {
		content := m.body[m.headLength:]
		if loc := reg.FindIndex(content); loc != nil && loc[0] == 0 {
			m.tag = string(content[:loc[1]])
		}
	}
	return
}

func (m *M) length() (l int, err error) {
	const funcTitle = packageTitle + "*M.init"
	if len(m.body) < 2 {
		err = errors.Wrap(errors.New("lenght part is two short"), funcTitle)
		return
	}
	l = int(binary.BigEndian.Uint16(m.body[:2]))
	return
}
