package message

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

type M struct {
	body       []byte
	headLength int
}

func New(headLength int, body []byte) *M {
	return &M{
		body:       body,
		headLength: headLength,
	}
}

func (m *M) Content() (content string, err error) {
	const funcTitle = packageTitle + "*M.Content"
	contentLength, err := m.length()
	if err != nil {
		err = errors.Wrapf(err, "%s: wrong message format", funcTitle)
		return
	}
	if contentLength == 0 {
		return
	}
	totalLength := m.headLength + contentLength
	if len(m.body) < totalLength {
		err = errors.Wrapf(errors.New("content is too short"), "%s: wrong message format", funcTitle)
		return
	}
	content = string(m.body[m.headLength:totalLength])
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
