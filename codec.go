package yamlvector

import (
	"bytes"
	"io"

	"github.com/koykov/vector"
)

const (
	flagEscapedString = 0
	flagFoldBlock     = 1
)

type Codec struct {
	vector.BaseCodec
}

func (Codec) Decode(p *vector.Byteptr) ([]byte, error) {
	return p.RawBytes(), nil
}

func (Codec) DecodeBool(p *vector.Byteptr) (bool, error) {
	b := p.RawBytes()
	return bytes.Equal(b, btrue) || bytes.Equal(b, bTrue) || bytes.Equal(b, bTRUE) ||
		bytes.Equal(b, bOn) || bytes.Equal(b, bon) || bytes.Equal(b, bON), nil
}

func (Codec) Beautify(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, true)
}

func (Codec) Marshal(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, false)
}
