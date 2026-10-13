package yamlvector

import (
	"io"

	"github.com/koykov/vector"
)

const (
	flagEscapedString = 0
	flagFoldBlock     = 1
)

type Helper struct{}

var (
	helper = Helper{}
)

func (h Helper) Indirect(p *vector.Byteptr) []byte {
	return p.RawBytes()
}

func (h Helper) Beautify(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, true)
}

func (h Helper) Marshal(w io.Writer, node *vector.Node) error {
	return serialize(w, node, 0, false)
}
